package chainstatev2

import (
	"github.com/hacash/core/fields"
)

// 写入包含交易哈希
func (cs *ChainState) ContainTxHash(txhx fields.Hash, blkhei fields.BlockHeight) error {
	query, e1 := cs.txhxchkDB.CreateNewQueryInstance(txhx)
	if e1 != nil {
		return e1 // error
	}
	defer query.Destroy()
	heibtx, _ := blkhei.Serialize()
	e4 := query.Save(heibtx)
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
	hei, e := cs.ReadTxBelongHeightByHash(txhx)
	if e != nil {
		return false, e
	}
	if hei > 0 {
		return true, nil
	}
	return false, nil // not find
}

// 检查交易所属区块高度
func (cs *ChainState) ReadTxBelongHeightByHash(txhx fields.Hash) (fields.BlockHeight, error) {

	query, e1 := cs.txhxchkDB.CreateNewQueryInstance(txhx)
	if e1 != nil {
		return 0, nil // error
	}
	defer query.Destroy()
	vdatas, e2 := query.Find()
	if e2 != nil {
		return 0, nil // error
	}
	if vdatas == nil {
		if cs.base != nil {
			return cs.base.ReadTxBelongHeightByHash(txhx) // check base
		} else {
			return 0, nil // not find
		}
	}
	if len(vdatas) == int(fields.BlockHeightSize) {
		var hei fields.BlockHeight
		hei.Parse(vdatas, 0)
		return hei, nil // find
	} else {
		return 0, nil // not find
	}
}

func (cs *ChainState) ReadTransactionBytesByHash(txhx fields.Hash) (fields.BlockHeight, []byte, error) {
	hei, bts, e := cs.BlockStore().ReadTransactionBytesByHash(txhx)
	return fields.BlockHeight(hei), bts, e
}
