package chainstatev3

import (
	"fmt"
	"github.com/hacash/chain/leveldb"
	"github.com/hacash/core/interfaces"
	"sync"
)

func (s *ChainState) TraversalCopy(src interfaces.ChainState) error {
	sta, ok := src.(*ChainState)
	if !ok {
		return fmt.Errorf("src chainstate is must *ChainState")
	}
	return s.TraversalCopyByObj(sta)
}

func (s *ChainState) TraversalCopyByObj(src *ChainState) error {

	if src.IsImmutable() {
		panic("TraversalCopy src state cannot must use LevelDB!")
	}

	myIsImm := s.IsImmutable()

	s.statusMux.Lock()
	defer s.statusMux.Unlock()

	if myIsImm {
		// leveldb
		e := s.traversalCopyMemToLevelUnsafe(s.ldb, src.memdb)
		if e != nil {
			return e
		}
	} else {
		// memdb
		src.memdb.Range(func(key, value interface{}) bool {
			s.memdb.Store(key, value)
			return true
		})
	}
	// err
	return nil
}

func (s ChainState) traversalCopyMemToLevelUnsafe(ldb *leveldb.DB, mem *sync.Map) error {
	var e error = nil
	batch := leveldb.MakeBatch(1)
	mem.Range(func(key, value interface{}) bool {
		// save to leveldb
		k := key.(string)
		v := value.(*MemoryStorageItem)
		//var e error = nil
		if v.IsDelete {
			// delete
			batch.Delete([]byte(k))
			//e = ldb.Delete([]byte(k), nil)
		} else {
			// save & update
			batch.Put([]byte(k), v.Value)
			//e = ldb.Put([]byte(k), v.Value, nil)
		}
		//if e != nil {
		//	return false
		//}
		return true
	})
	e = ldb.Write(batch, nil)
	// err
	return e
}
