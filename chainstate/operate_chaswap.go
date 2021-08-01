package chainstate

import (
	"github.com/hacash/core/fields"
	"github.com/hacash/core/stores"
)

//
func (cs *ChainState) Chaswap(chaswap_id fields.HashHalfChecker) *stores.Chaswap {
	query, e1 := cs.chaswapDB.CreateNewQueryInstance(chaswap_id)
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
			return cs.base.Chaswap(chaswap_id) // check base
		} else {
			return nil // not find
		}
	}
	if len(vdatas) == 0 {
		return nil // error
	}
	var stoitem stores.Chaswap
	_, e3 := stoitem.Parse(vdatas, 0)
	if e3 != nil {
		return nil // error
	}
	// return ok
	return &stoitem
}

//
func (cs *ChainState) ChaswapCreate(chaswap_id fields.HashHalfChecker, chaswap *stores.Chaswap) error {
	query, e1 := cs.chaswapDB.CreateNewQueryInstance(chaswap_id)
	if e1 != nil {
		return e1 // error
	}
	defer query.Destroy()
	stodatas, e3 := chaswap.Serialize()
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
func (cs *ChainState) ChaswapDelete(chaswap_id fields.HashHalfChecker) error {
	query, e1 := cs.chaswapDB.CreateNewQueryInstance(chaswap_id)
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
