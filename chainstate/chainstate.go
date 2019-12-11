package chainstate

import (
	"fmt"
	"github.com/hacash/chain/hashtreedb"
	"github.com/hacash/chain/tinykvdb"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/stores"
	"math/rand"
	"os"
	"path"
	"strconv"
)

type ChainState struct {
	temporaryDataDir string

	// parent state
	base *ChainState

	// config
	config *ChainStateConfig

	// status
	laststatusDB *tinykvdb.TinyKVDB

	// state
	balanceDB *hashtreedb.HashTreeDB
	diamondDB *hashtreedb.HashTreeDB
	channelDB *hashtreedb.HashTreeDB

	// store
	datastore interfaces.ChainStore

	// data hold
	pendingBlockHeight *uint64
	pendingBlockHash   fields.Hash

	submitStoreDiamond *stores.DiamondSmelt

	prev288BlockTimestamp   uint64
	lastestBlockHeadAndMeta interfaces.Block
	lastestDiamond          *stores.DiamondSmelt
}

func NewChainState(cnf *ChainStateConfig) (*ChainState, error) {
	return newChainStateEx(cnf, false)
}

func newChainStateEx(cnf *ChainStateConfig, isSubBranchTemporary bool) (*ChainState, error) {
	var laststatusDB *tinykvdb.TinyKVDB = nil
	// is temp state
	var temporaryDataDir string
	if isSubBranchTemporary {
		randstr := strconv.FormatUint(uint64(rand.Uint32()), 10)
		temporaryDataDir = path.Join(os.TempDir(), "hacash_state_temp_"+randstr)
		cnf.Datadir = temporaryDataDir
	} else {
		// laststatusDB
		lsdb, lserr := tinykvdb.NewTinyKVDB(path.Join(cnf.Datadir, "laststatus"))
		if lserr != nil {
			return nil, lserr
		}
		laststatusDB = lsdb
	}
	// balanceDB
	// fmt.Println("balanceDB dir:", path.Join(cnf.Datadir, "balance"))
	// os.MkdirAll(path.Join(cnf.Datadir, "balance"), os.ModePerm)
	blscnf := hashtreedb.NewHashTreeDBConfig(path.Join(cnf.Datadir, "balance"), stores.BalanceSize, 21)
	blscnf.KeyReverse = true
	if !isSubBranchTemporary {
		blscnf.FileDividePartitionLevel = 2
	} else {
		blscnf.ForbidGC = true
		blscnf.KeepDeleteMark = true
		blscnf.SaveMarkBeforeValue = true
	}
	balanceDB := hashtreedb.NewHashTreeDB(blscnf)
	// diamondDB
	dmdcnf := hashtreedb.NewHashTreeDBConfig(path.Join(cnf.Datadir, "diamond"), stores.DiamondSize, 6)
	dmdcnf.KeyPrefixSupplement = 10
	if !isSubBranchTemporary {
		dmdcnf.FileDividePartitionLevel = 1
	} else {
		dmdcnf.ForbidGC = true
		dmdcnf.KeepDeleteMark = true
		dmdcnf.SaveMarkBeforeValue = true
	}
	diamondDB := hashtreedb.NewHashTreeDB(dmdcnf)
	// channelDB
	chlcnf := hashtreedb.NewHashTreeDBConfig(path.Join(cnf.Datadir, "channel"), stores.ChannelSize, 16)
	if !isSubBranchTemporary {
		chlcnf.FileDividePartitionLevel = 1
	} else {
		chlcnf.ForbidGC = true
		chlcnf.KeepDeleteMark = true
		chlcnf.SaveMarkBeforeValue = true
	}
	channelDB := hashtreedb.NewHashTreeDB(chlcnf)
	// return ok
	cs := &ChainState{
		temporaryDataDir:      temporaryDataDir,
		base:                  nil,
		config:                cnf,
		laststatusDB:          laststatusDB,
		balanceDB:             balanceDB,
		diamondDB:             diamondDB,
		channelDB:             channelDB,
		prev288BlockTimestamp: 0,
		pendingBlockHeight:    nil,
	}
	return cs, nil
}

// destory temporary
func (cs *ChainState) DestoryTemporary() {
	if len(cs.temporaryDataDir) == 0 || cs.base == nil {
		return
	}
	// remove temp data dir
	e1 := os.RemoveAll(cs.temporaryDataDir)
	if e1 != nil {
		fmt.Println(e1)
	}
}

// chain data store
func (cs *ChainState) ChainStore() interfaces.ChainStore {
	if cs.datastore != nil {
		return cs.datastore
	}
	if cs.base != nil {
		cs.datastore = cs.base.ChainStore() // copy
		return cs.datastore
	}
	return cs.datastore
}

func (cs *ChainState) SetChainStore(store interfaces.ChainStore) error {
	if cs.base != nil {
		return fmt.Errorf("Can only be set chainstore in the final state.")
	}
	cs.datastore = store
	return nil
}

// merge write
func (cs *ChainState) MergeCoverWriteChainState(src *ChainState) error {
	// copy status
	if src.pendingBlockHeight != nil {
		cs.SetPendingBlockHeight(*src.pendingBlockHeight)
	}
	if src.pendingBlockHash != nil {
		cs.SetPendingBlockHash(src.pendingBlockHash)
	}
	if src.submitStoreDiamond != nil {
		e := cs.SetPendingSubmitStoreDiamond(src.submitStoreDiamond)
		if e != nil {
			return e
		}
	}
	if src.lastestBlockHeadAndMeta != nil {
		e := cs.SetLastestBlockHeadAndMeta(src.lastestBlockHeadAndMeta)
		if e != nil {
			return e
		}
	}
	if src.lastestDiamond != nil {
		e := cs.SetLastestDiamond(src.lastestDiamond)
		if e != nil {
			return e
		}
	}

	//  COPY COVER WRITE STATE

	e1 := cs.balanceDB.TraversalCopy(src.balanceDB, false)
	if e1 != nil {
		return e1
	}
	e2 := cs.diamondDB.TraversalCopy(src.diamondDB, false)
	if e2 != nil {
		return e2
	}
	e3 := cs.channelDB.TraversalCopy(src.channelDB, false)
	if e3 != nil {
		return e3
	}

	// copy ok

	return nil
}

// fork sub
func (cs *ChainState) NewSubBranchTemporaryChainState() (*ChainState, error) {

	tempcnf := NewChainStateConfig("")
	newTempState, err1 := newChainStateEx(tempcnf, true)
	if err1 != nil {
		return nil, err1
	}
	// set base
	newTempState.base = cs
	// ok
	return newTempState, nil
}

// submit to write disk
func (cs *ChainState) SubmitDataStoreWriteToInvariableDisk(block interfaces.Block) error {
	if cs.base != nil {
		return fmt.Errorf("Can only be saved in the final state.")
	}
	//
	store := cs.ChainStore()
	if store == nil {
		return fmt.Errorf("Not find ChainStore object.")
	}
	// save status
	e0 := cs.SetLastestBlockHeadAndMeta(block)
	if e0 != nil {
		return e0
	}
	e1 := cs.IncompleteSaveLastestBlockHeadAndMeta()
	if e1 != nil {
		return e1
	}
	e2 := cs.IncompleteSaveLastestDiamond()
	if e2 != nil {
		return e2
	}
	// save diamond
	if cs.submitStoreDiamond != nil {
		if cs.pendingBlockHash == nil {
			return fmt.Errorf("Block pending hash not set.")
		}
		cs.submitStoreDiamond.ContainBlockHash = cs.pendingBlockHash // copy
		e := store.SaveDiamond(cs.submitStoreDiamond)
		if e != nil {
			return e
		}
		cs.submitStoreDiamond = nil // reset
	}
	// save block data
	e3 := store.SaveBlockUniteTransactions(block)
	if e3 != nil {
		return e3
	}
	// clear
	cs.pendingBlockHash = nil
	cs.pendingBlockHeight = nil
	// ok
	return nil
}
