package chainstatev3

import (
	"fmt"
	"github.com/hacash/core/interfacev3"
)

// 保存在磁盘
func (s *ChainState) ImmutableWriteToDisk() (interfacev3.ChainStateImmutable, error) {
	if s.base.IsImmutable() == false {
		return nil, fmt.Errorf("State parent is not immutable.")
	}
	if s.IsImmutable() == true {
		return nil, fmt.Errorf("State self is immutable.")
	}

	s.statusMux.Lock()
	defer s.statusMux.Unlock()

	// copy save
	e := s.traversalCopyMemToLevelUnsafe(s.base.ldb, s.memdb)
	if e != nil {
		// err
		return nil, e
	}

	// update ptr
	s.ldb = s.base.ldb
	s.memdb = nil // delete

	// ok
	return s, nil
}
