package chainstate

import (
	"github.com/hacash/core/fields"
	"github.com/hacash/core/stores"
)

func (cs *ChainState) GetPendingBlockHeight() uint64 {
	if cs.pendingBlockHeight == nil {
		if cs.base != nil {
			return cs.base.GetPendingBlockHeight()
		}
		return 0
	}
	return *cs.pendingBlockHeight
}

func (cs *ChainState) SetPendingBlockHeight(height uint64) {
	setheight := height + 1 - 1 // copy
	cs.pendingBlockHeight = &setheight
}

////////////////////////////////////

func (cs *ChainState) GetPendingBlockHash() fields.Hash {
	if cs.pendingBlockHeight == nil {
		if cs.base != nil {
			return cs.base.GetPendingBlockHash()
		}
		return nil
	}
	return cs.pendingBlockHash
}

func (cs *ChainState) SetPendingBlockHash(hash fields.Hash) {
	cs.pendingBlockHash = hash
	// set diamond contail block hash
	cs.setPendingSubmitStoreDiamondContainBlockHash(hash)
}

func (cs *ChainState) setPendingSubmitStoreDiamondContainBlockHash(hash fields.Hash) {
	if cs.submitStoreDiamond != nil {
		cs.submitStoreDiamond.ContainBlockHash = hash
		return
	}
	if cs.base != nil {
		cs.base.setPendingSubmitStoreDiamondContainBlockHash(hash)
		return
	}
	// not exist diamond
	return
}

////////////////////////////////////

func (cs *ChainState) GetPendingSubmitStoreDiamond() (*stores.DiamondSmelt, error) {
	if cs.submitStoreDiamond == nil {
		if cs.base != nil {
			return cs.base.GetPendingSubmitStoreDiamond()
		}
		return nil, nil // not find
	}
	return cs.submitStoreDiamond, nil
}

func (cs *ChainState) SetPendingSubmitStoreDiamond(diamond *stores.DiamondSmelt) error {
	cs.submitStoreDiamond = diamond
	return nil
}
