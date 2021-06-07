package chainstate

import (
	"github.com/hacash/core/fields"
	"github.com/hacash/core/stores"
)

// DiamondLending 查询
func (cs *ChainState) UserLending(lendid fields.Bytes17) *stores.UserLending {
	query, e1 := cs.usrlendDB.CreateNewQueryInstance(lendid)
	if e1 != nil {
		return nil // error
	}
	defer query.Destroy()
	vdatas, e2 := query.Find()
	if e2 != nil {
		return nil // error
	}
	if vdatas == nil {
		if cs.base != nil {
			return cs.base.UserLending(lendid) // check base
		} else {
			return nil // not find
		}
	}
	var stoitem stores.UserLending
	_, e3 := stoitem.Parse(vdatas, 0)
	if e3 != nil {
		return nil // error
	}
	// return ok
	return &stoitem
}

// 创建 Diamond Lending
func (cs *ChainState) UserLendingCreate(lendid fields.Bytes17, stoitem *stores.UserLending) error {
	query, e1 := cs.usrlendDB.CreateNewQueryInstance(lendid)
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
func (cs *ChainState) UserLendingUpdate(lendid fields.Bytes17, stoitem *stores.UserLending) error {
	return cs.UserLendingCreate(lendid, stoitem)
}

// 删除
func (cs *ChainState) UserLendingDelete(lendid fields.Bytes17) error {
	query, e1 := cs.usrlendDB.CreateNewQueryInstance(lendid)
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
