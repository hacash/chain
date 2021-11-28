package chainstatev3

import (
	"github.com/hacash/core/fields"
	"github.com/hacash/core/stores"
)

func (cs *ChainState) Chaswap(chaswap_id fields.HashHalfChecker) (*stores.Chaswap, error) {
	value, ok, e := cs.find(KeySuffixType_chaswap, chaswap_id)
	if e != nil {
		return nil, e
	}
	if !ok {
		return nil, nil // not find
	}
	// parse
	var stoitem stores.Chaswap
	_, e = stoitem.Parse(value, 0)
	if e != nil {
		return nil, e // error
	}
	// return ok
	return &stoitem, nil
}

func (cs *ChainState) ChaswapCreate(chaswap_id fields.HashHalfChecker, chaswap *stores.Chaswap) error {
	stodatas, e := chaswap.Serialize()
	if e != nil {
		return e // error
	}
	// do save
	return cs.save(KeySuffixType_chaswap, chaswap_id, stodatas)
}

func (cs *ChainState) ChaswapUpdate(chaswap_id fields.HashHalfChecker, chaswap *stores.Chaswap) error {
	return cs.ChaswapCreate(chaswap_id, chaswap)
}

func (cs *ChainState) ChaswapDelete(chaswap_id fields.HashHalfChecker) error {
	return cs.delete(KeySuffixType_chaswap, chaswap_id)
}
