package chainstatev2

import (
	"fmt"
	"github.com/hacash/core/stores"
)

func (cs *ChainState) UpdateSetTotalSupply(totalobj *stores.TotalSupply) error {
	cs.totalSupply = totalobj
	cs.totalSupply_forSave = totalobj
	return nil
}

func (cs *ChainState) ReadTotalSupply() (*stores.TotalSupply, error) {
	if cs.totalSupply != nil {
		return cs.totalSupply, nil
	}
	if cs.base != nil {
		parentObj, e1 := cs.base.ReadTotalSupply()
		if e1 != nil {
			return nil, e1
		}
		if parentObj == nil {
			return nil, fmt.Errorf("not find TotalSupply from store.")
		}
		// Copy obj to avoid repeated addition to
		cs.totalSupply = parentObj.Clone()
		return cs.totalSupply, nil
	}
	// read from status db
	ins, e3 := cs.laststatusDB.CreateNewQueryInstance([]byte(LastestStatusKeyName_total_supply))
	if e3 != nil {
		return nil, e3
	}
	vdatas, e4 := ins.Find()
	if e4 != nil {
		return nil, e4
	}
	ttsupplyobj := stores.NewTotalSupplyStoreData()
	if vdatas == nil {
		// return genesis block
		return ttsupplyobj, nil
	}
	_, err1 := ttsupplyobj.Parse(vdatas, 0)
	if err1 != nil {
		return nil, err1
	}
	// cache set
	cs.totalSupply = ttsupplyobj
	return ttsupplyobj, nil
}

func (cs *ChainState) IncompleteSaveTotalSupply() error {
	if cs.laststatusDB == nil {
		return fmt.Errorf("cs.laststatusDB is not init.")
	}
	if cs.totalSupply_forSave == nil {
		return nil // not set
	}
	stodatas, e2 := cs.totalSupply_forSave.Serialize()
	if e2 != nil {
		return e2
	}
	cs.totalSupply_forSave = nil // clean data
	// save
	//fmt.Println("cs *ChainState) IncompleteSaveTotalSupply() error", stodatas)
	ins, e3 := cs.laststatusDB.CreateNewQueryInstance([]byte(LastestStatusKeyName_total_supply))
	if e3 != nil {
		return e3
	}
	e4 := ins.Save(stodatas)
	if e4 != nil {
		return e4
	}
	// ok
	return nil

}
