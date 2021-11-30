package chainstatev3

import (
	"github.com/hacash/core/fields"
)

// block data store
func (cs *ChainState) SaveMoveBTCBelongTxHash(trsno uint32, txhash []byte) error {

	num := fields.VarUint4(trsno)
	numbts, e := num.Serialize()
	if e != nil {
		return e // error
	}
	return cs.save(KeySuffixType_movebtc, numbts, txhash)
}

// block data store
func (cs *ChainState) ReadMoveBTCTxHashByTrsNo(trsno uint32) ([]byte, error) {

	num := fields.VarUint4(trsno)
	numbts, e := num.Serialize()
	if e != nil {
		return nil, e // error
	}

	value, ok, e := cs.find(KeySuffixType_movebtc, numbts)
	if e != nil {
		return nil, e
	}
	if !ok {
		return nil, nil // not find
	}
	return value, nil

}
