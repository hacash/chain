package chainstatev3

import (
	"github.com/hacash/core/fields"
	"github.com/hacash/core/stores"
)

func (cs *ChainState) Balance(addr fields.Address) (*stores.Balance, error) {
	value, ok, e := cs.find(KeySuffixType_balance, addr)
	if e != nil {
		return nil, e
	}
	if !ok {
		return stores.NewEmptyBalance(), nil // not find
	}

	var stoitem stores.Balance
	_, e = stoitem.Parse(value, 0)
	if e != nil {
		return nil, e // error
	}
	// return ok
	return &stoitem, nil
}

func (cs *ChainState) BalanceSet(addr fields.Address, balance *stores.Balance) error {
	stodatas, e := balance.Serialize()
	if e != nil {
		return e // error
	}
	// do save
	return cs.save(KeySuffixType_balance, addr, stodatas)
}

func (cs *ChainState) BalanceDel(addr fields.Address) error {
	return cs.delete(KeySuffixType_balance, addr)
}
