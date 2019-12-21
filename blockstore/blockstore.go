package blockstore

import (
	"github.com/hacash/chain/biglogdb"
	"github.com/hacash/chain/hashtreedb"
	"github.com/hacash/core/blocks"
	"github.com/hacash/core/stores"
	"path"
)

type BlockStore struct {

	// config
	config *BlockStoreConfig

	// data store
	blockdataDB  *biglogdb.BigLogDB
	trsdataptrDB *hashtreedb.HashTreeDB
	blknumhashDB *hashtreedb.HashTreeDB
	diamondDB    *hashtreedb.HashTreeDB
	diamondnumDB *hashtreedb.HashTreeDB
}

func NewBlockStore(cnf *BlockStoreConfig) (*BlockStore, error) {
	// create blockdataDB
	blcnf := biglogdb.NewBigLogDBConfig(path.Join(cnf.Datadir, "blockdata"), 32)
	blcnf.LogHeadMaxSize = blocks.BlockHeadSize
	blcnf.BlockPartFileMaxSize = 1024 * 1024 * 100 // 100MB
	blcnf.FileDividePartitionLevel = 1
	blcnf.KeyReverse = true // reverse key
	blockdataDB, e0 := biglogdb.NewBigLogDB(blcnf)
	if e0 != nil {
		return nil, e0
	}
	// create trsdataptrDB
	tdrcnf := hashtreedb.NewHashTreeDBConfig(path.Join(cnf.Datadir, "trsdataptr"), 5+biglogdb.LogFilePtrSeekSize, 32)
	tdrcnf.FileDividePartitionLevel = 2
	trsdataptrDB := hashtreedb.NewHashTreeDB(tdrcnf)
	// create blknumhashDB
	bnhcnf := hashtreedb.NewHashTreeDBConfig(path.Join(cnf.Datadir, "blocknum"), 32, 8)
	tdrcnf.KeyPrefixSupplement = 8
	blknumhashDB := hashtreedb.NewHashTreeDB(bnhcnf)
	// create diamondDB
	dmdcnf := hashtreedb.NewHashTreeDBConfig(path.Join(cnf.Datadir, "diamond"), stores.DiamondSmeltSize, 6)
	tdrcnf.KeyPrefixSupplement = 11
	diamondDB := hashtreedb.NewHashTreeDB(dmdcnf)
	// create diamondnumDB
	dmdnumcnf := hashtreedb.NewHashTreeDBConfig(path.Join(cnf.Datadir, "diamondnum"), 6, 4)
	tdrcnf.KeyPrefixSupplement = 4
	diamondnumDB := hashtreedb.NewHashTreeDB(dmdnumcnf)
	// return ok
	cs := &BlockStore{
		config:       cnf,
		blockdataDB:  blockdataDB,
		trsdataptrDB: trsdataptrDB,
		blknumhashDB: blknumhashDB,
		diamondDB:    diamondDB,
		diamondnumDB: diamondnumDB,
	}
	return cs, nil
}
