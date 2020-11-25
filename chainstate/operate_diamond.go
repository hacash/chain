package chainstate

import (
	"github.com/hacash/core/fields"
	"github.com/hacash/core/stores"
)

//
func (cs *ChainState) Diamond(diamond fields.Bytes6) *stores.Diamond {
	query, e1 := cs.diamondDB.CreateNewQueryInstance(diamond)
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
			return cs.base.Diamond(diamond) // check base
		} else {
			return nil // not find
		}
	}
	if len(vdatas) < stores.DiamondSize {
		return nil // error
	}
	var stoitem stores.Diamond
	_, e3 := stoitem.Parse(vdatas, 0)
	if e3 != nil {
		return nil // error
	}
	// return ok
	return &stoitem
}

//
func (cs *ChainState) DiamondSet(diamond_name fields.Bytes6, diamond *stores.Diamond) error {
	query, e1 := cs.diamondDB.CreateNewQueryInstance(diamond_name)
	if e1 != nil {
		return e1 // error
	}
	defer query.Destroy()
	stodatas, e3 := diamond.Serialize()
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

//
func (cs *ChainState) DiamondDel(diamond_name fields.Bytes6) error {
	query, e1 := cs.diamondDB.CreateNewQueryInstance(diamond_name)
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
