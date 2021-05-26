package biglogdb

import (
	"encoding/binary"
	"fmt"
	"github.com/hacash/chain/hashtreedb"
	"os"
	"path"
	"strconv"
	"sync"
)

/**
 * config
 */
type BigLogDBConfig struct {
	UseLevelDB               bool
	DataDir                  string
	KeySize                  uint8
	KeyReverse               bool
	LogHeadMaxSize           int
	BlockPartFileMaxSize     int64
	FileDividePartitionLevel uint8
}

func NewBigLogDBConfig(
	DataDir string,
	keySize uint8,
) *BigLogDBConfig {
	return &BigLogDBConfig{
		DataDir:                  DataDir,
		KeySize:                  keySize,
		KeyReverse:               false,
		LogHeadMaxSize:           0,
		FileDividePartitionLevel: 1,
		BlockPartFileMaxSize:     1024 * 1024 * 20, // 20MB
	}
}

type storefilecache struct {
	file *os.File
	num  uint32
}

/**
 * big log db
 */
type BigLogDB struct {
	config *BigLogDBConfig

	//////////////////////

	bashhashtreedb *hashtreedb.HashTreeDB

	headstat        *os.File
	headstatFileNum uint32

	storefilecache []storefilecache

	wlock sync.Mutex
}

// create DataBase
func NewBigLogDB(config *BigLogDBConfig) (*BigLogDB, error) {

	os.MkdirAll(config.DataDir, os.ModePerm)
	hsdbdir := path.Join(config.DataDir, "INDEXS")
	hsdbcnf := hashtreedb.NewHashTreeDBConfig(
		hsdbdir,
		LogFilePtrSeekSize+uint32(config.LogHeadMaxSize),
		config.KeySize,
	)
	// copy cnf con
	hsdbcnf.LevelDB = config.UseLevelDB
	hsdbcnf.KeyReverse = config.KeyReverse
	hsdbcnf.FileDividePartitionLevel = config.FileDividePartitionLevel
	// new tree db
	basedb := hashtreedb.NewHashTreeDB(hsdbcnf)
	db := &BigLogDB{
		config:          config,
		bashhashtreedb:  basedb,
		headstat:        nil,
		headstatFileNum: 4294967295, // not use
		storefilecache:  make([]storefilecache, 0),
	}
	return db, nil
}

func (db *BigLogDB) getStoreFileByNum(num uint32) (*os.File, error) {
	for i, v := range db.storefilecache {
		if v.num == num {
			if i > 0 && len(db.storefilecache) > 1 {
				db.storefilecache[i-1], db.storefilecache[i] =
					db.storefilecache[i], db.storefilecache[i-1]
				// change sort
			}
			return v.file, nil // get cache
		}
	}
	// open file
	ptfilen := path.Join(db.config.DataDir, "part"+strconv.Itoa(int(num))+".dat")
	f1, e1 := os.OpenFile(ptfilen, os.O_RDWR|os.O_CREATE, 0777)
	if e1 != nil {
		return nil, e1
	}
	sto := storefilecache{
		file: f1,
		num:  num,
	}
	db.storefilecache = append([]storefilecache{sto}, db.storefilecache...)
	if len(db.storefilecache) > 5 {
		db.storefilecache = db.storefilecache[0:5] // max size = 5
	}
	return f1, nil
}

////////////////////////////////////////////////////////

func (db *BigLogDB) getFileNumFile() (*os.File, error) {
	if db.headstat != nil {
		return db.headstat, nil
	}
	hdfilen := path.Join(db.config.DataDir, "HEAD.dat")
	f1, e1 := os.OpenFile(hdfilen, os.O_RDWR|os.O_CREATE, 0777)
	if e1 != nil {
		return nil, e1
	}
	db.headstat = f1
	return f1, nil
}

func (db *BigLogDB) GetFileNum() (uint32, error) {
	if db.headstatFileNum < 4294967295 {
		return db.headstatFileNum, nil
	}
	f1, e1 := db.getFileNumFile()
	if e1 != nil {
		return 0, e1
	}
	fstat, e0 := f1.Stat()
	if e0 != nil {
		return 0, e0
	}
	numdts := make([]byte, 4)
	if fstat.Size() == 0 {
		// init 0
		binary.BigEndian.PutUint32(numdts, 0)
		_, e3 := f1.Write(numdts)
		if e3 != nil {
			return 0, e3
		}
		db.headstatFileNum = 0
		return 0, nil

	} else {
		rn, e2 := f1.Read(numdts)
		if e2 != nil {
			return 0, e2
		}
		if rn != 4 && rn != 0 {
			return 0, fmt.Errorf("head file break.")
		}
		if rn == 4 {
			db.headstatFileNum = binary.BigEndian.Uint32(numdts)
			return db.headstatFileNum, nil
		}
	}
	return 0, fmt.Errorf("read file error.")
}

func (db *BigLogDB) SetFileNum(newnum uint32) error {
	f1, e1 := db.getFileNumFile()
	if e1 != nil {
		return e1
	}
	numdts := make([]byte, 4)
	binary.BigEndian.PutUint32(numdts, newnum)
	wn, e2 := f1.Write(numdts)
	if e2 != nil {
		return e2
	}
	if wn != 4 {
		return fmt.Errorf("head file break.")
	}
	db.headstatFileNum = newnum
	return nil
}

func (db *BigLogDB) Close() {
	if db.bashhashtreedb != nil {
		db.bashhashtreedb.Close()
	}
	if db.headstat != nil {
		db.headstat.Close()
	}
}
