package chainstate

import (
	"fmt"
	"github.com/hacash/core/blocks"
	"github.com/hacash/core/genesis"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/stores"
)

const (
	LastestStatusKeyName_lastest_block_head = "lastest_block_head"
	LastestStatusKeyName_lastest_diamond    = "lastest_diamond"
)

// status
func (cs *ChainState) SetLastestBlockHead(blockhead interfaces.Block) error {
	cs.lastestBlockHead = blockhead
	return nil
}

func (cs *ChainState) IncompleteSaveLastestBlockHead() error {
	if cs.laststatusDB == nil {
		return fmt.Errorf("cs.laststatusDB is not init.")
	}
	if cs.lastestBlockHead == nil {
		return nil // not set
	}
	stodatas, e2 := cs.lastestBlockHead.SerializeHead()
	if e2 != nil {
		return e2
	}
	// save
	e3 := cs.laststatusDB.Set([]byte(LastestStatusKeyName_lastest_block_head), stodatas)
	if e3 != nil {
		return e3
	}
	// ok
	return nil
}

func (cs *ChainState) ReadLastestBlockHead() (interfaces.Block, error) {
	if cs.lastestBlockHead != nil {
		return cs.lastestBlockHead, nil
	}
	if cs.base != nil {
		return cs.base.ReadLastestBlockHead()
	}
	// read from status db
	vdatas, e2 := cs.laststatusDB.Get([]byte(LastestStatusKeyName_lastest_block_head))
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
	tarblk, _, err := blocks.ParseBlockHead(vdatas, 0)
	if err != nil {
		return nil, err
	}
	// cache set
	cs.lastestBlockHead = tarblk
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
	if cs.lastestBlockHead != nil {
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
