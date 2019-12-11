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
)

func (cs *ChainState) SetLastestBlockHeadAndMeta(blockmeta interfaces.Block) error {
	cs.lastestBlockHeadAndMeta = blockmeta
	return nil
}

func (cs *ChainState) IncompleteSaveLastestBlockHeadAndMeta() error {
	if cs.laststatusDB == nil {
		return fmt.Errorf("cs.laststatusDB is not init.")
	}
	if cs.lastestBlockHeadAndMeta == nil {
		return nil // not set
	}
	stodatas, e2 := cs.lastestBlockHeadAndMeta.SerializeExcludeTransactions()
	if e2 != nil {
		return e2
	}
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
	tarblk, _, err1 := blocks.ParseBlockHead(vdatas, 0)
	if err1 != nil {
		return nil, err1
	}
	_, err1 = tarblk.ParseExcludeTransactions(vdatas, 0)
	if err1 != nil {
		return nil, err1
	}
	// cache set
	cs.lastestBlockHeadAndMeta = tarblk
	return tarblk, nil
}

/////////////////////////////////////////////////////////////////////////

func (cs *ChainState) SetLastestDiamond(diamond *stores.DiamondSmelt) error {
	cs.lastestDiamond = diamond
	return nil
}

func (cs *ChainState) IncompleteSaveLastestDiamond() error {
	if cs.laststatusDB == nil {
		return fmt.Errorf("cs.laststatusDB is not init.")
	}
	if cs.lastestDiamond == nil {
		return nil // not set
	}
	stodatas, e2 := cs.lastestDiamond.Serialize()
	if e2 != nil {
		return e2
	}
	// save
	e3 := cs.laststatusDB.Set([]byte(LastestStatusKeyName_lastest_diamond), stodatas)
	if e3 != nil {
		return e3
	}
	// ok
	return nil
}

func (cs *ChainState) ReadLastestDiamond() (*stores.DiamondSmelt, error) {
	if cs.lastestBlockHeadAndMeta != nil {
		return cs.lastestDiamond, nil
	}
	if cs.base != nil {
		return cs.base.ReadLastestDiamond()
	}
	// read from status db
	vdatas, e2 := cs.laststatusDB.Get([]byte(LastestStatusKeyName_lastest_diamond))
	if e2 != nil {
		return nil, e2
	}
	if vdatas == nil {
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
	cs.lastestDiamond = &diamond
	return &diamond, nil

}
