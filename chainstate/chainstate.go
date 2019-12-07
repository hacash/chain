package chainstate

import (
	"fmt"
	"github.com/hacash/chain/chainstore"
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
	datastore *chainstore.ChainStore

	// data hold
	pendingBlockHeight *uint64
	pendingBlockHash   fields.Hash

	submitStoreDiamond *stores.DiamondSmelt

	lastestBlockHead interfaces.Block
	lastestDiamond   *stores.DiamondSmelt
}

func NewChainState(cnf *ChainStateConfig) (*ChainState, error) {
	return newChainStateEx(cnf, false)
}

func newChainStateEx(cnf *ChainStateConfig, isSubBranchTemporary bool) (*ChainState, error) {
	// is temp state
	var temporaryDataDir string
	if isSubBranchTemporary {
		randstr := strconv.FormatUint(uint64(rand.Uint32()), 10)
		temporaryDataDir = path.Join(os.TempDir(), "hacash_state_temp_"+randstr)
		cnf.Datadir = temporaryDataDir
	}
	// laststatusDB
	laststatusDB, lserr := tinykvdb.NewTinyKVDB(path.Join(cnf.Datadir, "laststatus"))
	if lserr != nil {
		return nil, lserr
	}
	// balanceDB
	blscnf := hashtreedb.NewHashTreeDBConfig(path.Join(cnf.Datadir, "balance"), stores.BalanceSize, 21)
	blscnf.KeyReverse = true
	if !isSubBranchTemporary {
		blscnf.FileDividePartitionLevel = 2
	}
	balanceDB := hashtreedb.NewHashTreeDB(blscnf)
	// diamondDB
	dmdcnf := hashtreedb.NewHashTreeDBConfig(path.Join(cnf.Datadir, "diamond"), stores.DiamondSize, 6)
	dmdcnf.KeyPrefixSupplement = 10
	if !isSubBranchTemporary {
		dmdcnf.FileDividePartitionLevel = 1
	}
	diamondDB := hashtreedb.NewHashTreeDB(dmdcnf)
	// channelDB
	chlcnf := hashtreedb.NewHashTreeDBConfig(path.Join(cnf.Datadir, "channel"), stores.ChannelSize, 16)
	if !isSubBranchTemporary {
		chlcnf.FileDividePartitionLevel = 1
	}
	channelDB := hashtreedb.NewHashTreeDB(chlcnf)
	// return ok
	cs := &ChainState{
		temporaryDataDir:   temporaryDataDir,
		base:               nil,
		config:             cnf,
		laststatusDB:       laststatusDB,
		balanceDB:          balanceDB,
		diamondDB:          diamondDB,
		channelDB:          channelDB,
		pendingBlockHeight: nil,
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
func (cs *ChainState) GetChainStore() *chainstore.ChainStore {
	if cs.datastore != nil {
		return cs.datastore
	}
	if cs.base != nil {
		cs.datastore = cs.base.GetChainStore() // copy
		return cs.datastore
	}
	return cs.datastore
}

func (cs *ChainState) SetChainStore(store *chainstore.ChainStore) error {
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
		cs.SetPendingSubmitStoreDiamond(src.submitStoreDiamond)
	}
	if src.lastestBlockHead != nil {
		e := cs.SetLastestBlockHead(src.lastestBlockHead)
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

	//  TODO: COPY COVER WRITE STATE

	return nil
}

// fork sub
func (cs *ChainState) NewSubBranchTemporaryChainState(abspath string) (*ChainState, error) {

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
	store := cs.GetChainStore()
	if store == nil {
		return fmt.Errorf("Not find ChainStore object.")
	}
	// save status
	e1 := cs.IncompleteSaveLastestBlockHead()
	if e1 != nil {
		return e1
	}
	e2 := cs.IncompleteSaveLastestDiamond()
	if e2 != nil {
		return e2
	}
	// save block data
	e3 := store.SaveBlockUniteTransactions(block)
	if e3 != nil {
		return e3
	}
	// ok
	return nil
}
