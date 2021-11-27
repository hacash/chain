package chainstatev2

import (
	"fmt"
	"github.com/hacash/core/fields"
)

// 写入包含交易哈希
func (cs *ChainState) ContainTxHash(txhx fields.Hash) error {
	query, e1 := cs.txhxchkDB.CreateNewQueryInstance(txhx)
	if e1 != nil {
		return e1 // error
	}
	defer query.Destroy()
	e4 := query.Save([]byte{1})
	if e4 != nil {
		return e4 // error
	}
	// ok
	return nil
}

// 移除交易
func (cs *ChainState) RemoveTxHash(txhx fields.Hash) error {
	query, e1 := cs.txhxchkDB.CreateNewQueryInstance(txhx)
	if e1 != nil {
		return e1 // error
	}
	defer query.Destroy()
	e2 := query.Delete()
	if e2 != nil {
		return e2 // error
	}
	return nil
}

// 检查交易是否上链
func (cs *ChainState) CheckTxHash(txhx fields.Hash) (bool, error) {
	query, e1 := cs.txhxchkDB.CreateNewQueryInstance(txhx)
	if e1 != nil {
		return false, nil // error
	}
	defer query.Destroy()
	vdatas, e2 := query.Find()
	if e2 != nil {
		return false, nil // error
	}
	if vdatas == nil {
		if cs.base != nil {
			return cs.base.CheckTxHash(txhx) // check base
		} else {
			return false, nil // not find
		}
	}
	if len(vdatas) != 1 {
		return false, fmt.Errorf("vdatas len error")
	}
	if vdatas[0] == 0 {
		return false, nil
	} else {
		return true, nil
	}
}
