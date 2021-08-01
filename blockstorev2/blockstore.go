package blockstorev2

import (
	"github.com/hacash/chain/biglogdb"
	"github.com/hacash/chain/leveldb"
	"github.com/hacash/chain/statedomaindb"
	"github.com/hacash/core/blocks"
	"path"
)

type BlockStore struct {

	// config
	config *BlockStoreConfig

	// data store
	blockdataDB *biglogdb.BigLogDB

	// store
	trsdataptrDB *statedomaindb.StateDomainDB
	blknumhashDB *statedomaindb.StateDomainDB
	diamondDB    *statedomaindb.StateDomainDB
	diamondnumDB *statedomaindb.StateDomainDB

	// btc move log
	btcmovelogDB        *statedomaindb.StateDomainDB
	btcmovelogTotalPage int // 最大数据页码
}

func NewBlockStoreOfBlockDataDB(basedir string, ldb *leveldb.DB) (*biglogdb.BigLogDB, error) {
	var blockPartFileMaxSize int64 = 1024 * 1024 * 100 // 100MB
	blcnf := biglogdb.NewBigLogDBConfig(path.Join(basedir, "blockdata"), 32, blockPartFileMaxSize)
	blcnf.LogHeadMaxSize = blocks.BlockHeadSize
	return biglogdb.NewBigLogDBByLevelDB(blcnf, "blockdata", ldb)
}

func NewBlockStore(cnf *BlockStoreConfig) (*BlockStore, error) {

	// create leveldb
	useldb, e0 := leveldb.OpenFile(cnf.Datadir, nil)
	if e0 != nil {
		return nil, e0
	}

	// create blockdataDB
	blockdataDB, e0 := NewBlockStoreOfBlockDataDB(cnf.Datadir, useldb)
	if e0 != nil {
		return nil, e0
	}

	// trsdataptrDB
	tdrcnf := statedomaindb.NewStateDomainDBConfig("trsdataptr", 5+biglogdb.LogFilePtrSeekSize, 32)
	tdrcnf.LevelDB = true
	trsdataptrDB := statedomaindb.NewStateDomainDB(tdrcnf, useldb)

	// blocknumDB
	bnhcnf := statedomaindb.NewStateDomainDBConfig("blocknum", 32, 8)
	bnhcnf.LevelDB = true
	blknumhashDB := statedomaindb.NewStateDomainDB(bnhcnf, useldb)

	// diamondDB
	dmdcnf := statedomaindb.NewStateDomainDBConfig("diamond", 0, 6)
	dmdcnf.LevelDB = true
	diamondDB := statedomaindb.NewStateDomainDB(dmdcnf, useldb)

	// diamondnumDB
	dmdnumcnf := statedomaindb.NewStateDomainDBConfig("diamondnum", 0, 4)
	dmdnumcnf.LevelDB = true
	diamondnumDB := statedomaindb.NewStateDomainDB(dmdnumcnf, useldb)

	// btcmovelogsDB
	btmvcnf := statedomaindb.NewStateDomainDBConfig("btcmovelog", 0, 0)
	btmvcnf.LevelDB = true
	btcmovelogDB := statedomaindb.NewStateDomainDB(btmvcnf, useldb)

	// return ok
	cs := &BlockStore{
		config:              cnf,
		blockdataDB:         blockdataDB,
		trsdataptrDB:        trsdataptrDB,
		blknumhashDB:        blknumhashDB,
		diamondDB:           diamondDB,
		diamondnumDB:        diamondnumDB,
		btcmovelogDB:        btcmovelogDB,
		btcmovelogTotalPage: -1, //
	}
	return cs, nil
}

// 创建一个用于更新数据库版本的区块存储器
func NewBlockStoreForUpdateDatabaseVersion(cnf *BlockStoreConfig) (*BlockStore, error) {

	// create leveldb
	useldb, e0 := leveldb.OpenFile(cnf.Datadir, nil)
	if e0 != nil {
		return nil, e0
	}

	// create blockdataDB
	blockdataDB, e0 := NewBlockStoreOfBlockDataDB(cnf.Datadir, useldb)
	if e0 != nil {
		return nil, e0
	}

	bnhcnf := statedomaindb.NewStateDomainDBConfig("blocknum", 32, 8)
	bnhcnf.LevelDB = true
	blknumhashDB := statedomaindb.NewStateDomainDB(bnhcnf, useldb)

	// return ok
	cs := &BlockStore{
		config:              cnf,
		blockdataDB:         blockdataDB,
		blknumhashDB:        blknumhashDB,
		btcmovelogTotalPage: -1, //
	}
	return cs, nil
}

func (cs *BlockStore) Close() {
	if cs.blockdataDB != nil {
		cs.blockdataDB.Close()
	}
	if cs.trsdataptrDB != nil {
		cs.trsdataptrDB.Close()
	}
	if cs.blknumhashDB != nil {
		cs.blknumhashDB.Close()
	}
	if cs.diamondDB != nil {
		cs.diamondDB.Close()
	}
	if cs.diamondnumDB != nil {
		cs.diamondnumDB.Close()
	}
	if cs.btcmovelogDB != nil {
		cs.btcmovelogDB.Close()
	}
}
