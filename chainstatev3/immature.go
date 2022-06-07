package chainstatev3

import (
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
)

// Query the specified status tree location
func (s *ChainState) SearchBaseStateByBlockHashObj(hx fields.Hash) (*ChainState, error) {
	if s.pending != nil {
		if s.pending.GetPendingBlockHash().Equal(hx) {
			// State is itself
			return s, nil
		}
	}
	// Query sub status
	for _, sta := range s.childs {
		ptr, e := sta.SearchBaseStateByBlockHashObj(hx)
		if e != nil {
			return nil, e
		}
		if ptr != nil {
			return ptr, nil // Found
		}
	}
	// not found
	return nil, nil
}

func (s *ChainState) SearchBaseStateByBlockHash(hx fields.Hash) (interfaces.ChainState, error) {
	obj, e := s.SearchBaseStateByBlockHashObj(hx)
	return obj, e
}

// Traversing immature block hash
func (s *ChainState) SeekImmatureBlockHashs() ([]fields.Hash, error) {
	s.statusMux.RLock()
	defer s.statusMux.RUnlock()

	var hxs = make([]fields.Hash, 0)
	e := s.doSeekImmatureBlockHashs(&hxs) // 递归遍历
	// De duplication and de invalidation
	hxmaps := make(map[string]bool)
	newhxlist := make([]fields.Hash, 0)
	for _, v := range hxs {
		if len(v) == fields.HashSize && v.NotZeroBlank() {
			if _, ok := hxmaps[string(v)]; !ok {
				newhxlist = append(newhxlist, v)
			}
			hxmaps[string(v)] = true
		}
	}
	// return
	return newhxlist, e
}

func (s *ChainState) doSeekImmatureBlockHashs(hxs *[]fields.Hash) error {
	for _, child := range s.childs {
		hx := child.pending.GetPendingBlockHash()
		if len(hx) != fields.HashSize {
			continue // Invalid hash
		}
		*hxs = append(*hxs, hx)
		e := child.doSeekImmatureBlockHashs(hxs)
		if e != nil {
			return e
		}
	}
	return nil
}
