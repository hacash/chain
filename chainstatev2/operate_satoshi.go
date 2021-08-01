package chainstatev2

/*

//
func (cs *ChainState) Satoshi(addr fields.Address) *stores.Satoshi {
	query, e1 := cs.satoshiDB.CreateNewQueryInstance(addr)
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
			return cs.base.Satoshi(addr) // check base
		} else {
			return stores.NewEmptySatoshi() // not find
		}
	}
	if len(vdatas) < stores.SatoshiSize {
		return nil // error
	}
	var stoitem stores.Satoshi
	_, e3 := stoitem.Parse(vdatas, 0)
	if e3 != nil {
		return nil // error
	}
	// return ok
	return &stoitem
}

//
func (cs *ChainState) SatoshiSet(addr fields.Address, satoshi *stores.Satoshi) error {
	query, e1 := cs.satoshiDB.CreateNewQueryInstance(addr)
	if e1 != nil {
		return e1 // error
	}
	defer query.Destroy()
	stodatas, e3 := satoshi.Serialize()
	if e3 != nil {
		return e3 // error
	}
	_, e4 := query.Save(stodatas)
	if e4 != nil {
		return e4 // error
	}
	// ok
	return nil
}

//
func (cs *ChainState) SatoshiDel(addr fields.Address) error {
	query, e1 := cs.satoshiDB.CreateNewQueryInstance(addr)
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

*/
