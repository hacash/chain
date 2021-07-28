package chainstate

import (
	"github.com/hacash/core/fields"
	"github.com/hacash/core/stores"
)

//
func (cs *ChainState) Balance(addr fields.Address) *stores.Balance {
	query, e1 := cs.balanceDB.CreateNewQueryInstance(addr)
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
			return cs.base.Balance(addr) // check base
		} else {
			return stores.NewEmptyBalance() // not find
		}
	}
	if len(vdatas) == 0 {
		return nil // error
	}
	var stoitem stores.Balance
	_, e3 := stoitem.Parse(vdatas, 0)
	if e3 != nil {
		return nil // error
	}
	// return ok
	return &stoitem
}

//
func (cs *ChainState) BalanceSet(addr fields.Address, balance *stores.Balance) error {
	query, e1 := cs.balanceDB.CreateNewQueryInstance(addr)
	if e1 != nil {
		return e1 // error
	}
	defer query.Destroy()
	stodatas, e3 := balance.Serialize()
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
func (cs *ChainState) BalanceDel(addr fields.Address) error {
	query, e1 := cs.balanceDB.CreateNewQueryInstance(addr)
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
