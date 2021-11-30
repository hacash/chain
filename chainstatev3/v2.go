package chainstatev3

import (
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/stores"
)

func (cs *ChainState) GetPendingBlockHeight() uint64 {
	return cs.GetPending().GetPendingBlockHeight()
}

func (cs *ChainState) GetPendingBlockHash() fields.Hash {
	return cs.GetPending().GetPendingBlockHash()
}

func (cs *ChainState) ReadLastestBlockHeadMetaForRead() (interfaces.BlockHeadMetaRead, error) {

	return cs.GetPending().GetPendingBlockHead().(interfaces.BlockHeadMetaRead), nil
}

func (cs *ChainState) ReadLastestDiamond() (*stores.DiamondSmelt, error) {
	last, e := cs.LatestStatusRead()
	if e != nil {
		return nil, e
	}
	return last.ReadLastestDiamond(), nil
}
