package chainstatev3

import (
	"fmt"
	"github.com/hacash/core/interfaces"
)

// 保存在磁盘
func (s *ChainState) ImmutableWriteToDisk() (interfaces.ChainStateImmutable, error) {
	obj, e := s.ImmutableWriteToDiskObj()
	return obj, e
}

// 保存在磁盘
func (s *ChainState) ImmutableWriteToDiskObj() (*ChainState, error) {
	if s.base != nil && s.base.GetPendingBlockHeight() > 0 && s.base.IsImmutable() == false {
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
