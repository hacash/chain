package chainstatev3

import (
	"github.com/hacash/core/fields"
	"github.com/hacash/core/stores"
)

// Diamondlending query
func (cs *ChainState) BitcoinSystemLending(lendid fields.BitcoinSyslendId) (*stores.BitcoinSystemLending, error) {
	value, ok, e := cs.find(KeySuffixType_btclend, lendid)
	if e != nil {
		return nil, e
	}
	if !ok {
		return nil, nil // not find
	}
	// parse
	var stoitem stores.BitcoinSystemLending
	_, e = stoitem.Parse(value, 0)
	if e != nil {
		return nil, e // error
	}
	// return ok
	return &stoitem, nil
}

// Create diamond lending
func (cs *ChainState) BitcoinLendingCreate(lendid fields.BitcoinSyslendId, stoitem *stores.BitcoinSystemLending) error {
	stodatas, e := stoitem.Serialize()
	if e != nil {
		return e // error
	}
	// do save
	return cs.save(KeySuffixType_btclend, lendid, stodatas)
}

// to update
func (cs *ChainState) BitcoinLendingUpdate(lendid fields.BitcoinSyslendId, stoitem *stores.BitcoinSystemLending) error {
	return cs.BitcoinLendingCreate(lendid, stoitem)
}

// delete
func (cs *ChainState) BitcoinLendingDelete(lendid fields.BitcoinSyslendId) error {
	return cs.delete(KeySuffixType_btclend, lendid)
}
