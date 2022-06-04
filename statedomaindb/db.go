package statedomaindb

import (
	"fmt"
	"github.com/hacash/chain/leveldb"
)

type StateDomainDBConfig struct {
	// MemoryStorage
	MemoryStorage bool // Save data in memory
	// LevelDB
	LevelDB bool // Save data using level dB
	// size
	KeySize                  uint8  // Key value length < = 32
	SupplementalMaxValueSize uint32 // Data content length

	KeyDomainName string // Key prefix
}

func NewStateDomainDBConfig(
	keyDomainName string,
	mustMinValueSize uint32, // Data length that must be supplemented
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
	// In memory database
	if config.MemoryStorage {
		db.MemoryStorageDB = NewMemoryStorageDB()
		return db
	}
	// Using level dB
	if config.LevelDB {
		if ldb == nil {
			panic("Must give param ldb *leveldb.DB")
		}
		return db
	}

	panic("NewStateDomainDB  must use MemoryStorage or LevelDB!")

	// File database, data length
	// db.freshRecordDataSize()
	return db
}

// Create execution unit
func (db *StateDomainDB) CreateNewQueryInstance(key []byte) (*QueryInstance, error) {
	kz := int(db.config.KeySize)
	if kz > 0 && len(key) != kz {
		return nil, fmt.Errorf("len(domainkey)<%d> not more than db.config.KeySize<%d>", len(key), int(db.config.KeySize))
	}
	return newQueryInstance(db, key)
}

// close
func (db *StateDomainDB) Close() error {
	// Turn off leveldb
	if db.LevelDB != nil {
		db.LevelDB.Close()
	}
	return nil
}
