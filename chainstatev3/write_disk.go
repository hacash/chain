package chainstatev3

import (
	"fmt"
)

// 保存在磁盘
func (s *ChainState) ImmutableWriteToDisk() error {
	if s.base.IsImmutable() == false {
		return fmt.Errorf("State parent is not immutable.")
	}
	if s.IsImmutable() == true {
		return fmt.Errorf("State self is immutable.")
	}

	s.statusMux.Lock()
	defer s.statusMux.Unlock()

	// copy save
	e := s.traversalCopyMemToLevelUnsafe(s.base.ldb, s.memdb)
	if e != nil {
		// err
		return e
	}

	// update ptr
	s.ldb = s.base.ldb
	s.memdb = nil // delete

	// ok
	return nil
}
