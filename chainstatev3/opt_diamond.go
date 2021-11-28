package chainstatev3

import (
	"github.com/hacash/core/fields"
	"github.com/hacash/core/stores"
)

func (cs *ChainState) Diamond(diamond fields.DiamondName) (*stores.Diamond, error) {
	value, ok, e := cs.find(KeySuffixType_diamond, diamond)
	if e != nil {
		return nil, e
	}
	if !ok {
		return nil, nil // not find
	}
	// parse
	var stoitem stores.Diamond
	_, e = stoitem.Parse(value, 0)
	if e != nil {
		return nil, e // error
	}
	// return ok
	return &stoitem, nil
}

func (cs *ChainState) DiamondSet(diamond_name fields.DiamondName, diamond *stores.Diamond) error {
	stodatas, e := diamond.Serialize()
	if e != nil {
		return e // error
	}
	// do save
	return cs.save(KeySuffixType_diamond, diamond_name, stodatas)
}

func (cs *ChainState) DiamondDel(diamond_name fields.DiamondName) error {
	return cs.delete(KeySuffixType_diamond, diamond_name)
}
