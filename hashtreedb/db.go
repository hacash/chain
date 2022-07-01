package hashtreedb

import (
	"fmt"
	"github.com/hacash/chain/leveldb"
	"sync"
)

// A single file size supports at least 256^4 × five × About 80GB when 8 menuwide=8

/*
const (
	IndexItemSize int = 1 + 4              // Fixed
	IndexMenuSize int = 16 * IndexItemSize // Fixed
)

const (
	IndexItemTypeNull        = byte(0)
	IndexItemTypeBranch      = byte(1)
	IndexItemTypeValue       = byte(2)
	IndexItemTypeValueDelete = byte(3)
)

*/

type HashTreeDBConfig struct {
	// MemoryStorage
	MemoryStorage bool // Save data in memory
	// LevelDB
	LevelDB bool // Save data using level dB
	// size
	KeySize                  uint8  // Key value length < = 32
	SupplementalMaxValueSize uint32 // Data content length
	// key config
	//KeyReverse          bool  // key值倒序
	//KeyPrefixSupplement uint8 // key值前缀增补
	// opt config
	//SaveMarkBeforeValue bool // 储存原始的key值到Value前面 // 用于遍历改写
	//SaveKeyBeforeValue    bool // 储存原始的key值到Value前面 // 用于遍历改写
	//KeepDeleteMark            bool   // 删除也会保存key标记
	//TargetFilePackagePoolSize uint32 // 操作单例的缓存池大小

	// file config
	//FileDividePartitionLevel uint8  // 文件分区层级 0为不分区
	FileAbsPath string // File storage path
	//FileName                 string // 保存文件的名称
	// gc
	//ForbidGC bool // 禁止垃圾空间回收管理

	// other
	//hashSize         uint8
	//segmentValueSize uint32
}

func NewHashTreeDBConfig(
	fileAbsPath string,
	mustMinValueSize uint32, // Data length that must be supplemented
	keySize uint8,
) *HashTreeDBConfig {
	return &HashTreeDBConfig{
		FileAbsPath:              fileAbsPath,
		SupplementalMaxValueSize: mustMinValueSize,
		KeySize:                  keySize,
		//ForbidGC:                  false,
		//SaveMarkBeforeValue:       false,
		//TargetFilePackagePoolSize: 1,
		//KeyReverse:                false,
		//KeyPrefixSupplement:       0,
		//FileDividePartitionLevel:  0,
		//FileName:                  "blk",
		//KeepDeleteMark:            false,
	}
}

type HashTreeDB struct {
	config *HashTreeDBConfig // config

	// db in memory
	MemoryStorageDB *MemoryStorageDB

	// db in memory
	LevelDB          *leveldb.DB
	levelDBCreateMux sync.Mutex

	// file opt
	//filesOptLock   sync.Mutex
	//filesWriteLock sync.Map // map[string]*lockFilePkgItem

	//fileWriteLockCount  sync.Map // map[string]int         // 写文件锁数量统计
	//fileWriteLockMutexs sync.Map // map[string]*sync.Mutex // 写文件锁
	//targetFilePackageCache *TargetFilePackage // map[string]*TargetFilePackage // 暂时版本先只储存一个

	//existsFileKeys sync.Map // 已经存在的

	//HashSize   uint32 // 哈希大小 16,32,64,128,256
	//KeyReverse bool   // key值倒序
	//
	//SupplementalMaxValueSize uint32 // 最大数据尺寸大小 + hash32
	//
	//MenuWide uint8 // 单层索引宽度数（不可超过256）
	//
	//FilePartitionLevel uint32 // 文件分区层级 0为不分区
	//
	//FileAbsPath string // 文件的储存路径
	//FileName    string // 保存文件的名称
	//FileSuffix  string // 保存文件后缀名 .idx
	//
	//DeleteMark bool // 删除也会保存key标记

	////gc *GarbageCollectionDB
	//OpenGc       bool                            // 是否开启gc
	//gcPool       map[string]*GarbageCollectionDB // gc管理器
	//MaxNumGCPool int
	//
	//// fileLock
	//FileLock sync.Map // map[string]*sync.Mutex
}

// create DataBase
func NewHashTreeDB(config *HashTreeDBConfig) *HashTreeDB {
	db := &HashTreeDB{
		config: config,
	}
	// In memory database
	if config.MemoryStorage {
		db.MemoryStorageDB = NewMemoryStorageDB()
		return db
	}
	// Using level dB
	if config.LevelDB {
		//fmt.Println("config.LevelDB file path: ", config.FileAbsPath)
		// Create on demand when using
		return db
	}

	panic("NewHashTreeDB  must use LevelDB!")

	// File database, data length
	// db.freshRecordDataSize()
	return db
}

// Get or create a level dB object
func (db *HashTreeDB) GetOrCreateLevelDBwithPanic() *leveldb.DB {
	if db.LevelDB != nil {
		return db.LevelDB
	}
	db.levelDBCreateMux.Lock()
	defer db.levelDBCreateMux.Unlock()
	if db.LevelDB != nil {
		return db.LevelDB
	}
	ldb, err := leveldb.OpenFile(db.config.FileAbsPath, nil)
	if err != nil {
		fmt.Println("NewHashTreeDB leveldb.OpenFile Error:", err.Error())
		panic(err)
	}
	db.LevelDB = ldb
	return db.LevelDB
}

// Create execution unit
func (db *HashTreeDB) CreateNewQueryInstance(key []byte) (*QueryInstance, error) {
	if len(key) != int(db.config.KeySize) {
		return nil, fmt.Errorf("len(key)<%d> not more than db.config.KeySize<%d>", len(key), int(db.config.KeySize))
	}
	return newQueryInstance(db, key)
}

// fresh size config
/*
func (db *HashTreeDB) freshRecordDataSize() {
	if int(db.config.KeyPrefixSupplement)+int(db.config.KeySize) > 32 {
		panic("KeyPrefixSupplement + KeySize not more than 32.")
	}
	db.config.hashSize = db.config.KeyPrefixSupplement + db.config.KeySize
	// markSize? + KeySize + SupplementalMaxValueSize
	db.config.segmentValueSize = 0
	if db.config.SaveMarkBeforeValue {
		db.config.segmentValueSize += uint32(1)
	}
	db.config.segmentValueSize += uint32(db.config.KeySize) + db.config.SupplementalMaxValueSize
}
*/
// close
func (db *HashTreeDB) Close() error {
	// Turn off leveldb
	if db.LevelDB != nil {
		db.LevelDB.Close()
	}
	/* 其他
	db.filesWriteLock.Range(func(key, value interface{}) bool {
		var item = value.(*lockFilePkgItem)
		item.lock.Lock()
		item.count = 0
		item.targetFilePackageCache.Destroy()
		item.lock.Unlock()
		return true
	})
	db.filesWriteLock = sync.Map{}
	*/

	//if db.targetFilePackageCache != nil {
	//	db.targetFilePackageCache.Destroy() // close cache
	//	db.targetFilePackageCache = nil
	//}
	return nil
}
