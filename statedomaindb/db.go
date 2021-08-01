package statedomaindb

import (
	"fmt"
	"github.com/hacash/chain/leveldb"
)

type StateDomainDBConfig struct {
	// MemoryStorage
	MemoryStorage bool // 在内存内保存数据
	// LevelDB
	LevelDB bool // 使用 level db 保存数据
	// size
	KeySize                  uint8  // key值长度  <= 32
	SupplementalMaxValueSize uint32 // 数据内容长度

	KeyDomainName string // key 前缀
}

func NewStateDomainDBConfig(
	keyDomainName string,
	mustMinValueSize uint32, // 必须补足的数据长度
	keySize uint8,
) *StateDomainDBConfig {
	return &StateDomainDBConfig{
		KeyDomainName:            keyDomainName,
		SupplementalMaxValueSize: mustMinValueSize,
		KeySize:                  keySize,
	}
}

type StateDomainDB struct {
	config *StateDomainDBConfig // config

	// db in memory
	MemoryStorageDB *MemoryStorageDB

	// db in memory
	LevelDB *leveldb.DB
}

// create DataBase
func NewStateDomainDB(config *StateDomainDBConfig, ldb *leveldb.DB) *StateDomainDB {
	db := &StateDomainDB{
		config:  config,
		LevelDB: ldb,
	}
	// 内存数据库
	if config.MemoryStorage {
		db.MemoryStorageDB = NewMemoryStorageDB()
		return db
	}
	// 使用 level db
	if config.LevelDB {
		if ldb == nil {
			panic("Must give param ldb *leveldb.DB")
		}
		return db
	}

	panic("NewStateDomainDB  must use MemoryStorage or LevelDB!")

	// 文件数据库，数据长度
	// db.freshRecordDataSize()
	return db
}

// 创建执行单元
func (db *StateDomainDB) CreateNewQueryInstance(key []byte) (*QueryInstance, error) {
	kz := int(db.config.KeySize)
	if kz > 0 && len(key) != kz {
		return nil, fmt.Errorf("len(domainkey)<%d> not more than db.config.KeySize<%d>", len(key), int(db.config.KeySize))
	}
	return newQueryInstance(db, key)
}

// close
func (db *StateDomainDB) Close() error {
	// 关闭 leveldb
	if db.LevelDB != nil {
		db.LevelDB.Close()
	}
	return nil
}
