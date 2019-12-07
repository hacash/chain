package chainstore

import (
	"github.com/hacash/chain/biglogdb"
	"github.com/hacash/chain/hashtreedb"
	"github.com/hacash/core/blocks"
	"github.com/hacash/core/stores"
	"path"
)

type ChainStore struct {

	// config
	config *ChainStoreConfig

	// data store
	blockdataDB  *biglogdb.BigLogDB
	trsdataptrDB *hashtreedb.HashTreeDB
	blknumhashDB *hashtreedb.HashTreeDB
	diamondDB    *hashtreedb.HashTreeDB
	diamondnumDB *hashtreedb.HashTreeDB
}

func NewChainStore(cnf *ChainStoreConfig) (*ChainStore, error) {
	// create blockdataDB
	blcnf := biglogdb.NewBigLogDBConfig(path.Join(cnf.absdir, "blockdata"), 32)
	blcnf.LogHeadMaxSize = blocks.BlockHeadSize
	blcnf.BlockPartFileMaxSize = 1024 * 1024 * 100 // 100MB
	blcnf.FileDividePartitionLevel = 1
	blcnf.KeyReverse = true // reverse key
	blockdataDB, e0 := biglogdb.NewBigLogDB(blcnf)
	if e0 != nil {
		return nil, e0
	}
	// create trsdataptrDB
	tdrcnf := hashtreedb.NewHashTreeDBConfig(path.Join(cnf.absdir, "trsdataptr"), biglogdb.LogFilePtrSeekSize, 32)
	tdrcnf.FileDividePartitionLevel = 2
	trsdataptrDB := hashtreedb.NewHashTreeDB(tdrcnf)
	// create blknumhashDB
	bnhcnf := hashtreedb.NewHashTreeDBConfig(path.Join(cnf.absdir, "blocknum"), 32, 8)
	tdrcnf.KeyPrefixSupplement = 8
	blknumhashDB := hashtreedb.NewHashTreeDB(bnhcnf)
	// create diamondDB
	dmdcnf := hashtreedb.NewHashTreeDBConfig(path.Join(cnf.absdir, "diamond"), stores.DiamondSmeltSize, 6)
	tdrcnf.KeyPrefixSupplement = 11
	diamondDB := hashtreedb.NewHashTreeDB(dmdcnf)
	// create diamondnumDB
	dmdnumcnf := hashtreedb.NewHashTreeDBConfig(path.Join(cnf.absdir, "diamondnum"), 6, 4)
	tdrcnf.KeyPrefixSupplement = 4
	diamondnumDB := hashtreedb.NewHashTreeDB(dmdnumcnf)
	// return ok
	cs := &ChainStore{
		config:       cnf,
		blockdataDB:  blockdataDB,
		trsdataptrDB: trsdataptrDB,
		blknumhashDB: blknumhashDB,
		diamondDB:    diamondDB,
		diamondnumDB: diamondnumDB,
	}
	return cs, nil
}
