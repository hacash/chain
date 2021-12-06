package chainstatev3

import (
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
)

// 查询指定的状态树位置
func (s *ChainState) SearchBaseStateByBlockHashObj(hx fields.Hash) (*ChainState, error) {
	if s.pending != nil {
		if s.pending.GetPendingBlockHash().Equal(hx) {
			// 状态就是自身
			return s, nil
		}
	}
	// 查询子状态
	for _, sta := range s.childs {
		ptr, e := sta.SearchBaseStateByBlockHashObj(hx)
		if e != nil {
			return nil, e
		}
		if ptr != nil {
			return ptr, nil // 已经找到
		}
	}
	// 未找到
	return nil, nil
}

func (s *ChainState) SearchBaseStateByBlockHash(hx fields.Hash) (interfaces.ChainState, error) {
	obj, e := s.SearchBaseStateByBlockHashObj(hx)
	return obj, e
}

// 遍历不成熟的区块哈希
func (s *ChainState) SeekImmatureBlockHashs() ([]fields.Hash, error) {
	s.statusMux.RLock()
	defer s.statusMux.RUnlock()

	var hxs = make([]fields.Hash, 0)
	e := s.doSeekImmatureBlockHashs(&hxs) // 递归遍历
	// 去重去无效
	hxmaps := make(map[string]bool)
	newhxlist := make([]fields.Hash, 0)
	for _, v := range hxs {
		if len(v) == fields.HashSize {
			if _, ok := hxmaps[string(v)]; !ok {
				newhxlist = append(newhxlist, v)
			}
			hxmaps[string(v)] = true
		}
	}
	// 返回
	return newhxlist, e
}

func (s *ChainState) doSeekImmatureBlockHashs(hxs *[]fields.Hash) error {
	for _, child := range s.childs {
		hx := child.pending.GetPendingBlockHash()
		if len(hx) != fields.HashSize {
			continue // 无效的哈希
		}
		*hxs = append(*hxs, hx)
		e := child.doSeekImmatureBlockHashs(hxs)
		if e != nil {
			return e
		}
	}
	return nil
}
