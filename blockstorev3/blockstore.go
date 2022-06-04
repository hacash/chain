package blockstorev3

import (
	"fmt"
	leveldb "github.com/hacash/chain/leveldb"
	"sync"
)

type BlockStore struct {

	// config
	config *BlockStoreConfig

	// level db
	ldb *leveldb.DB

	btcmovelogTotalPage int // Maximum data page number

	statusMux *sync.RWMutex
}

func NewBlockStore(cnf *BlockStoreConfig) (*BlockStore, error) {
	store := &BlockStore{
		config:              cnf,
		ldb:                 nil,
		btcmovelogTotalPage: -1,
		statusMux:           &sync.RWMutex{},
	}

	useldb, e := leveldb.OpenFile(cnf.Datadir, nil)
	if e != nil {
		return nil, e
	}

	store.statusMux.Lock()
	store.ldb = useldb
	store.statusMux.Unlock()

	// ok
	return store, nil
}

func (bs *BlockStore) getDB() (*leveldb.DB, error) {

	bs.statusMux.RLock()
	var ldb = bs.ldb
	bs.statusMux.RUnlock()

	if ldb == nil {
		return nil, fmt.Errorf("level db has been closed or not init.")
	}

	return ldb, nil
}

func (bs *BlockStore) Close() {
	bs.statusMux.Lock()
	defer bs.statusMux.Unlock()

	if bs.ldb != nil {
		bs.ldb.Close()
		bs.ldb = nil
	}
}
