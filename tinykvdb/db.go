package tinykvdb

import (
	"fmt"
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
	abspath          string
	UseLevelDB       bool
	ldb              *leveldb.DB
	levelDBCreateMux sync.Mutex

	///////////
	bashhashtreedb *hashtreedb.HashTreeDB
	storefile      *os.File
	wlock          sync.Mutex
}

// create DataBase
func NewTinyKVDB(abspath string, UseLevelDB bool) (*TinyKVDB, error) {

	if UseLevelDB {

		// 返回
		return &TinyKVDB{
			abspath:    abspath,
			UseLevelDB: true,
			ldb:        nil, // 按需创建
		}, nil
	}

	panic("must use level db!")

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

// 获取或创建 level db 对象
func (db *TinyKVDB) GetOrCreateLevelDBwithPanic() *leveldb.DB {
	if db.ldb != nil {
		return db.ldb
	}
	db.levelDBCreateMux.Lock()
	defer db.levelDBCreateMux.Unlock()
	if db.ldb != nil {
		return db.ldb
	}
	leveldbobj, err := leveldb.OpenFile(db.abspath, nil)
	if err != nil {
		fmt.Println("NewTinyKVDB leveldb.OpenFile Error", err)
		panic(err)
	}
	db.ldb = leveldbobj
	return db.ldb
}

func (kv *TinyKVDB) Close() error {
	if kv.storefile != nil {
		kv.storefile.Close()
	}
	if kv.ldb != nil {
		kv.ldb.Close()
	}
	if kv.bashhashtreedb != nil {
		kv.bashhashtreedb.Close()
	}
	return nil
}
