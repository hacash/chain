package chainstate

import (
	"fmt"
	"github.com/hacash/core/blocks"
	"github.com/hacash/core/genesis"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/stores"
)

const (
	LastestStatusKeyName_lastest_block_head_meta = "lastest_block_head_meta"
	LastestStatusKeyName_lastest_diamond         = "lastest_diamond"
	LastestStatusKeyName_total_supply            = "total_supply"
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

/////////////////////////////////////////////////////////////////////////

func (cs *ChainState) SetLastestBlockHeadAndMeta(blockmeta interfaces.Block) error {
	cs.lastestBlockHeadAndMeta = blockmeta
	cs.lastestBlockHeadAndMeta_forSave = blockmeta
	return nil
}

func (cs *ChainState) IncompleteSaveLastestBlockHeadAndMeta() error {
	if cs.laststatusDB == nil {
		return fmt.Errorf("cs.laststatusDB is not init.")
	}
	if cs.lastestBlockHeadAndMeta_forSave == nil {
		return nil // not set
	}
	stodatas, e2 := cs.lastestBlockHeadAndMeta_forSave.SerializeExcludeTransactions()
	if e2 != nil {
		return e2
	}
	cs.lastestBlockHeadAndMeta_forSave = nil // clean data
	// save
	e3 := cs.laststatusDB.Set([]byte(LastestStatusKeyName_lastest_block_head_meta), stodatas)
	if e3 != nil {
		return e3
	}
	// ok
	return nil
}

func (cs *ChainState) ReadLastestBlockHeadAndMeta() (interfaces.Block, error) {
	if cs.lastestBlockHeadAndMeta != nil {
		return cs.lastestBlockHeadAndMeta, nil
	}
	if cs.base != nil {
		return cs.base.ReadLastestBlockHeadAndMeta()
	}
	// read from status db
	vdatas, e2 := cs.laststatusDB.Get([]byte(LastestStatusKeyName_lastest_block_head_meta))
	if e2 != nil {
		return nil, e2
	}
	if vdatas == nil {
		// return genesis block
		return genesis.GetGenesisBlock(), nil
	}
	if len(vdatas) < blocks.BlockHeadSize {
		return nil, fmt.Errorf("lastest_block_head store file error.")
	}
	tarblk, _, err1 := blocks.ParseExcludeTransactions(vdatas, 0)
	if err1 != nil {
		return nil, err1
	}
	// cache set
	cs.lastestBlockHeadAndMeta = tarblk
	return tarblk, nil
}

/////////////////////////////////////////////////////////////////////////

func (cs *ChainState) SetLastestDiamond(diamond *stores.DiamondSmelt) error {
	//fmt.Println("<<<<<<<<<<<<<   SetLastestDiamond")
	//fmt.Println("diamond", string(diamond.Diamond), diamond.PrevContainBlockHash.ToHex(), diamond.ContainBlockHash.ToHex())
	cs.lastestDiamond = diamond
	cs.lastestDiamond_forSave = diamond
	return nil
}

func (cs *ChainState) IncompleteSaveLastestDiamond() error {
	if cs.laststatusDB == nil {
		panic("cs.laststatusDB is not init.")
	}
	if cs.lastestDiamond_forSave == nil {
		return nil // not set
	}
	//fmt.Println("IncompleteSaveLastestDiamond  cs.pendingBlockHash ", cs.pendingBlockHash.ToHex())
	if cs.lastestDiamond_forSave.ContainBlockHash == nil {
		if cs.pendingBlockHash == nil {
			//return fmt.Errorf("Block pending hash not set.")
			panic("pending block hash not be set.")
		}
		cs.lastestDiamond_forSave.ContainBlockHash = cs.pendingBlockHash // copy
	}
	stodatas, e2 := cs.lastestDiamond_forSave.Serialize()
	if e2 != nil {
		return e2
	}
	cs.lastestDiamond_forSave = nil
	// save
	//fmt.Println("IncompleteSaveLastestDiamond", LastestStatusKeyName_lastest_diamond, stodatas)
	e3 := cs.laststatusDB.Set([]byte(LastestStatusKeyName_lastest_diamond), stodatas)
	if e3 != nil {
		return e3
	}
	// ok
	return nil
}

func (cs *ChainState) ReadLastestDiamond() (*stores.DiamondSmelt, error) {
	//fmt.Println("ReadLastestDiamond >>>>>>>")
	if cs.lastestDiamond != nil {
		return cs.lastestDiamond, nil
	}
	if cs.base != nil {
		return cs.base.ReadLastestDiamond()
	}
	// read from status db
	//fmt.Println("ReadLastestDiamond >>>>>>>  read from status db")
	vdatas, e2 := cs.laststatusDB.Get([]byte(LastestStatusKeyName_lastest_diamond))
	if e2 != nil {
		return nil, e2
	}
	if vdatas == nil {
		//fmt.Println("ReadLastestDiamond   return nil, nil // first one")
		return nil, nil // first one
	}
	if len(vdatas) < stores.DiamondSmeltSize {
		return nil, fmt.Errorf("lastest_diamond store file error.")
	}
	var diamond stores.DiamondSmelt
	_, err := diamond.Parse(vdatas, 0)
	if err != nil {
		return nil, err
	}
	// cache set
	//fmt.Println("ReadLastestDiamond ", diamond)
	cs.lastestDiamond = &diamond
	return &diamond, nil

}
