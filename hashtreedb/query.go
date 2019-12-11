package hashtreedb

import (
	"sync"
)

// 查询实例

type QueryInstance struct {
	db *HashTreeDB

	key        []byte
	hash       []byte
	fileKey    string
	filePath   string
	searchHash []byte

	targetFilePackage         *TargetFilePackage
	targetFileWriteJustUnlock *sync.Mutex

	// search cache
	searchResultCache *FindValueOffsetItem
}

/**
 * QueryInstance {
 *     Find()
 *     Save(valuebytes []byte)
 * }
 */
func newQueryInstance(db *HashTreeDB, key []byte) (*QueryInstance, error) {

	ins := &QueryInstance{
		db:                db,
		key:               key,
		searchResultCache: nil,
	}
	ins.hash = db.convertKeyToHash(key)
	ins.filePath, ins.fileKey, ins.searchHash = db.locateTargetFilePath(ins.hash)
	//fmt.Println("newQueryInstance searchHash ", ins.searchHash)
	// 等待获取文件控制
	lock, err := db.waitForTakeControlOfFile(ins)
	if err != nil {
		return nil, err
	}
	ins.targetFileWriteJustUnlock = lock
	// 返回使用
	return ins, nil
}

// 关闭
func (ins *QueryInstance) Destroy() {
	// 释放文件控制
	ins.db.releaseControlOfFile(ins)
	// 清空数据
	ins.db = nil
	ins.key = nil
	ins.hash = nil
	ins.fileKey = ""
	ins.filePath = ""
	ins.searchHash = nil
	ins.targetFilePackage = nil
	ins.targetFileWriteJustUnlock = nil
	ins.searchResultCache = nil
}
