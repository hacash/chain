package chainstatev3

import (
	"bytes"
	"fmt"
	"github.com/hacash/chain/leveldb"
)

const (
	KeySuffixType_balance = "balance"
	KeySuffixType_diamond = "diamond"
	KeySuffixType_channel = "channel"
	KeySuffixType_movebtc = "movebtc"
	KeySuffixType_lockbls = "lockbls"
	KeySuffixType_dmdlend = "dmdlend"
	KeySuffixType_btclend = "btclend"
	KeySuffixType_usrlend = "usrlend"
	KeySuffixType_chaswap = "chaswap"
	KeySuffixType_txhxchk = "txhxchk"
)

/*********************************************/

type MemoryStorageItem struct {
	IsDelete bool   // 已经删除标记
	Value    []byte // 数据
}

func NewDeleteMarkItem() *MemoryStorageItem {
	return &MemoryStorageItem{
		IsDelete: true,
		Value:    nil,
	}
}
func NewStorageItem(value []byte) *MemoryStorageItem {
	return &MemoryStorageItem{
		IsDelete: false,
		Value:    value,
	}
}

/*********************************************/

func keyfix(k []byte, suffix string) string {
	buf := bytes.NewBuffer(k)
	buf.Write([]byte(suffix))
	return string(buf.Bytes())
}

/*********************************************/

/**
 * find save update delete
 */
func (s *ChainState) save(suffix string, k, value []byte) error {
	usekey := keyfix(k, suffix)

	s.statusMux.RLock()
	var mem = s.memdb
	var ldb = s.ldb
	s.statusMux.RUnlock()

	if mem != nil {
		// save
		mem.Store(usekey, NewStorageItem(value))
		return nil
	}

	if ldb != nil {
		// delete from leveldb
		return ldb.Put([]byte(usekey), value, nil)
	}

	// error
	return fmt.Errorf("State has not memdb or leveldb both.")
}

func (s *ChainState) delete(suffix string, k []byte) error {
	usekey := keyfix(k, suffix)

	s.statusMux.RLock()
	var mem = s.memdb
	var ldb = s.ldb
	s.statusMux.RUnlock()

	if mem != nil {
		// add delete mark
		mem.Store(usekey, NewDeleteMarkItem())
		return nil
	}

	if ldb != nil {
		// delete from leveldb
		return ldb.Delete([]byte(usekey), nil)
	}

	// error
	return fmt.Errorf("State has not memdb or leveldb both.")
}

func (s ChainState) find(suffix string, k []byte) ([]byte, bool, error) {
	usekey := keyfix(k, suffix)

	s.statusMux.RLock()
	var mem = s.memdb
	var ldb = s.ldb
	var base = s.base
	s.statusMux.RUnlock()

	// check memdb
	if mem != nil {
		value, ok := mem.Load(usekey)
		if ok {
			vobj := value.(*MemoryStorageItem)
			if vobj.IsDelete {
				return nil, false, nil
			} else {
				return vobj.Value, true, nil
			}
		}
	}
	// check leveldb
	if ldb != nil {
		value, e := ldb.Get([]byte(usekey), nil)
		if e != nil {
			if e == leveldb.ErrNotFound {
				return nil, false, nil // not find
			} else {
				return nil, false, e // read error
			}
		}
		// find
		return value, true, nil
	}

	// find in parent
	if base != nil {
		return base.find(suffix, k)
	}

	// not find
	return nil, false, nil
}
