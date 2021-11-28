package chainstatev3

import (
	"github.com/hacash/core/fields"
	"github.com/hacash/core/stores"
)

func (cs *ChainState) DiamondSystemLending(dmdid fields.DiamondSyslendId) (*stores.DiamondSystemLending, error) {
	value, ok, e := cs.find(KeySuffixType_dmdlend, dmdid)
	if e != nil {
		return nil, e
	}
	if !ok {
		return nil, nil // not find
	}
	// parse
	var stoitem stores.DiamondSystemLending
	_, e = stoitem.Parse(value, 0)
	if e != nil {
		return nil, e // error
	}
	// return ok
	return &stoitem, nil
}

func (cs *ChainState) DiamondLendingCreate(dmdid fields.DiamondSyslendId, stoitem *stores.DiamondSystemLending) error {
	stodatas, e := stoitem.Serialize()
	if e != nil {
		return e // error
	}
	// do save
	return cs.save(KeySuffixType_dmdlend, dmdid, stodatas)
}

func (cs *ChainState) DiamondLendingUpdate(dmdid fields.DiamondSyslendId, stoitem *stores.DiamondSystemLending) error {
	return cs.DiamondLendingCreate(dmdid, stoitem)
}

func (cs *ChainState) DiamondLendingDelete(dmdid fields.DiamondSyslendId) error {
	return cs.delete(KeySuffixType_dmdlend, dmdid)
}
