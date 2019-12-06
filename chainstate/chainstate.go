package chainstate

import (
	"github.com/hacash/chain/chainstore"
	"github.com/hacash/chain/hashtreedb"
	"github.com/hacash/chain/tinykvdb"
)

type ChainState struct {

	// parent state
	base *ChainState

	// config
	config *ChainStateConfig

	// status
	laststatusDB *tinykvdb.TinyKVDB

	// state
	balanceDB *hashtreedb.HashTreeDB
	diamoneDB *hashtreedb.HashTreeDB
	channelDB *hashtreedb.HashTreeDB

	// store
	datastore *chainstore.ChainStore
}

func NewChainState(cnf *ChainStateConfig) *ChainState {

	cs := &ChainState{
		base:   nil,
		config: cnf,
	}
	return cs
}

// chain data store
func (cs *ChainState) GetChainStore() *chainstore.ChainStore {

	return nil
}

func (cs *ChainState) SetChainStore(store *chainstore.ChainStore) {

}

// merge write
func (cs *ChainState) MergeCoverWriteChainState(src *ChainState) error {

	return nil
}

// fork sub
func (cs *ChainState) NewSubBranchTemporaryChainState(abspath string) *ChainState {

	return nil
}
