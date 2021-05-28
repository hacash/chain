package chainstate

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
		// copy obj, 避免重复添加到
		cs.totalSupply = parentObj.Clone()
		return cs.totalSupply, nil
	}
	// read from status db
	vdatas, e2 := cs.laststatusDB.Get([]byte(LastestStatusKeyName_total_supply))
	if e2 != nil {
		return nil, e2
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
	e3 := cs.laststatusDB.Set([]byte(LastestStatusKeyName_total_supply), stodatas)
	if e3 != nil {
		return e3
	}
	// ok
	return nil

}
