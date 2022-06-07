package chainstatev2

import (
	"github.com/hacash/core/fields"
	"github.com/hacash/core/stores"
)

// Diamondlending query
func (cs *ChainState) DiamondSystemLending(dmdid fields.DiamondSyslendId) (*stores.DiamondSystemLending, error) {
	query, e1 := cs.dmdlendDB.CreateNewQueryInstance(dmdid)
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
			return cs.base.DiamondSystemLending(dmdid) // check base
		} else {
			return nil, nil // not find
		}
	}
	var stoitem stores.DiamondSystemLending
	_, e3 := stoitem.Parse(vdatas, 0)
	if e3 != nil {
		return nil, nil // error
	}
	// return ok
	return &stoitem, nil
}

// Create diamond lending
func (cs *ChainState) DiamondLendingCreate(dmdid fields.DiamondSyslendId, stoitem *stores.DiamondSystemLending) error {
	query, e1 := cs.dmdlendDB.CreateNewQueryInstance(dmdid)
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

// to update
func (cs *ChainState) DiamondLendingUpdate(dmdid fields.DiamondSyslendId, stoitem *stores.DiamondSystemLending) error {
	return cs.DiamondLendingCreate(dmdid, stoitem)
}

// delete
func (cs *ChainState) DiamondLendingDelete(dmdid fields.DiamondSyslendId) error {
	query, e1 := cs.dmdlendDB.CreateNewQueryInstance(dmdid)
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
