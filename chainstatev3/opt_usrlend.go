package chainstatev3

import (
	"github.com/hacash/core/fields"
	"github.com/hacash/core/stores"
)

// DiamondLending 查询
func (cs *ChainState) UserLending(lendid fields.UserLendingId) (*stores.UserLending, error) {
	value, ok, e := cs.find(KeySuffixType_usrlend, lendid)
	if e != nil {
		return nil, e
	}
	if !ok {
		return nil, nil // not find
	}
	// parse
	var stoitem stores.UserLending
	_, e = stoitem.Parse(value, 0)
	if e != nil {
		return nil, e // error
	}
	// return ok
	return &stoitem, nil
}

// 创建 Diamond Lending
func (cs *ChainState) UserLendingCreate(lendid fields.UserLendingId, stoitem *stores.UserLending) error {
	stodatas, e := stoitem.Serialize()
	if e != nil {
		return e // error
	}
	// do save
	return cs.save(KeySuffixType_usrlend, lendid, stodatas)
}

// 更新
func (cs *ChainState) UserLendingUpdate(lendid fields.UserLendingId, stoitem *stores.UserLending) error {
	return cs.UserLendingCreate(lendid, stoitem)
}

// 删除
func (cs *ChainState) UserLendingDelete(lendid fields.UserLendingId) error {
	return cs.delete(KeySuffixType_usrlend, lendid)
}
