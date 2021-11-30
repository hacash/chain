package chainstatev3

import (
	"github.com/hacash/core/stores"
)

func (cs *ChainState) UpdateSetTotalSupply(totalobj *stores.TotalSupply) error {
	databts, e := totalobj.Serialize()
	if e != nil {
		return e
	}
	return cs.save(KeySuffixType_totalsupply, []byte{1}, databts)
}

func (cs *ChainState) ReadTotalSupply() (*stores.TotalSupply, error) {

	datas, ok, e := cs.find(KeySuffixType_totalsupply, []byte{1})
	if e != nil {
		return nil, e
	}
	supplyobj := stores.NewTotalSupplyStoreData()
	if !ok {
		return supplyobj, nil
	}

	_, e = supplyobj.Parse(datas, 0)
	if e != nil {
		return nil, e
	}

	return supplyobj, nil
}
