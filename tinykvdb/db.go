package tinykvdb

import (
	"github.com/hacash/chain/hashtreedb"
	"github.com/hacash/chain/leveldb"
	"os"
	"strings"
	"sync"
)

const (
	ItemDelMark = uint8(3)
)

/**
 * small kv db
 */
type TinyKVDB struct {
	UseLevelDB bool
	ldb        *leveldb.DB

	///////////
	bashhashtreedb *hashtreedb.HashTreeDB
	storefile      *os.File
	wlock          sync.Mutex
}

// create DataBase
func NewTinyKVDB(abspath string, UseLevelDB bool) (*TinyKVDB, error) {

	if UseLevelDB {
		ldb, err := leveldb.OpenFile(abspath, nil)
		if err != nil {
			return nil, err
		}
		// 返回
		return &TinyKVDB{
			UseLevelDB: true,
			ldb:        ldb,
		}, nil
	}

	// create dir file
	os.MkdirAll(abspath, os.ModePerm)
	storefile, e1 := os.OpenFile(strings.TrimRight(abspath, "/")+"/storevalue.dat", os.O_RDWR|os.O_CREATE, 0777)
	if e1 != nil {
		return nil, e1
	}
	// hash tree db
	hxdbcnf := hashtreedb.NewHashTreeDBConfig(abspath, 1+4+4, 16)
	//hxdbcnf.FileDividePartitionLevel = 1
	bashhashtreedb := hashtreedb.NewHashTreeDB(hxdbcnf)
	// KVDB
	db := &TinyKVDB{
		UseLevelDB:     false,
		storefile:      storefile,
		bashhashtreedb: bashhashtreedb,
	}
	return db, nil
}

func (kv *TinyKVDB) Close() error {
	if kv.bashhashtreedb != nil {
		return kv.bashhashtreedb.Close()
	}
	return nil
}
