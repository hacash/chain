package tinykvdb

import (
	"github.com/hacash/chain/hashtreedb"
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
	bashhashtreedb *hashtreedb.HashTreeDB

	storefile *os.File

	wlock sync.Mutex
}

// create DataBase
func NewTinyKVDB(abspath string) (*TinyKVDB, error) {
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
		storefile:      storefile,
		bashhashtreedb: bashhashtreedb,
	}
	return db, nil
}


func (kv *TinyKVDB) Close() error {
	return kv.bashhashtreedb.Close()
}