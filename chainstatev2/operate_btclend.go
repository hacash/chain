package chainstatev2

import (
	"github.com/hacash/core/fields"
	"github.com/hacash/core/stores"
)

// DiamondLending 查询
func (cs *ChainState) BitcoinSystemLending(lendid fields.BitcoinSyslendId) (*stores.BitcoinSystemLending, error) {
	query, e1 := cs.btclendDB.CreateNewQueryInstance(lendid)
	if e1 != nil {
		return nil, nil // error
	}
	defer query.Destroy()
	vdatas, e2 := query.Find()
	if e2 != nil {
		return nil, nil // error
	}
	if vdatas == nil {
		if cs.base != nil {
			return cs.base.BitcoinSystemLending(lendid) // check base
		} else {
			return nil, nil // not find
		}
	}
	var stoitem stores.BitcoinSystemLending
	_, e3 := stoitem.Parse(vdatas, 0)
	if e3 != nil {
		return nil, nil // error
	}
	// return ok
	return &stoitem, nil
}

// 创建 Diamond Lending
func (cs *ChainState) BitcoinLendingCreate(lendid fields.BitcoinSyslendId, stoitem *stores.BitcoinSystemLending) error {
	query, e1 := cs.btclendDB.CreateNewQueryInstance(lendid)
	if e1 != nil {
		return e1 // error
	}
	defer query.Destroy()
	stodatas, e3 := stoitem.Serialize()
	if e3 != nil {
		return e3 // error
	}
	e4 := query.Save(stodatas)
	if e4 != nil {
		return e4 // error
	}
	// ok
	return nil
}

// 更新
func (cs *ChainState) BitcoinLendingUpdate(lendid fields.BitcoinSyslendId, stoitem *stores.BitcoinSystemLending) error {
	return cs.BitcoinLendingCreate(lendid, stoitem)
}

// 删除
func (cs *ChainState) BitcoinLendingDelete(lendid fields.BitcoinSyslendId) error {
	query, e1 := cs.btclendDB.CreateNewQueryInstance(lendid)
	if e1 != nil {
		return e1 // error
	}
	defer query.Destroy()
	e2 := query.Delete()
	if e2 != nil {
		return e2 // error
	}
	// ok
	return nil
}
