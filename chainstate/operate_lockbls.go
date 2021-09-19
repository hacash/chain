package chainstate

import (
	"github.com/hacash/core/fields"
	"github.com/hacash/core/stores"
)

// lockbls 查询
func (cs *ChainState) Lockbls(lkid fields.LockblsId) *stores.Lockbls {
	query, e1 := cs.lockblsDB.CreateNewQueryInstance(lkid)
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
			return cs.base.Lockbls(lkid) // check base
		} else {
			return nil // not find
		}
	}
	if len(vdatas) == 0 {
		return nil // error
	}
	var stoitem stores.Lockbls
	_, e3 := stoitem.Parse(vdatas, 0)
	if e3 != nil {
		return nil // error
	}
	// return ok
	return &stoitem
}

// 创建线性锁仓
func (cs *ChainState) LockblsCreate(lkid fields.LockblsId, stoitem *stores.Lockbls) error {
	query, e1 := cs.lockblsDB.CreateNewQueryInstance(lkid)
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

// 释放（取出部分任意可取额度）
func (cs *ChainState) LockblsUpdate(lkid fields.LockblsId, stoitem *stores.Lockbls) error {
	return cs.LockblsCreate(lkid, stoitem)
}

// 删除
func (cs *ChainState) LockblsDelete(lkid fields.LockblsId) error {
	query, e1 := cs.lockblsDB.CreateNewQueryInstance(lkid)
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
