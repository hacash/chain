package hashtreedb

// 查询实例

type QueryInstance struct {
	db *HashTreeDB

	key        []byte
	hash       []byte
	fileKey    string
	filePath   string
	searchHash []byte

	targetFilePackage *TargetFilePackage
	targetFileItem    *lockFilePkgItem

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
	// 如果是内存数据库，则不打开本地文件
	if db.config.MemoryStorage {
		return ins, nil
	}
	ins.hash = db.convertKeyToHash(key)
	ins.filePath, ins.fileKey, ins.searchHash = db.locateTargetFilePath(ins.hash)
	//fmt.Println("newQueryInstance searchHash ", ins.searchHash)
	// 等待获取文件控制
	fileitem, err := db.waitForTakeControlOfFile(ins)
	if err != nil {
		return nil, err
	}
	ins.targetFileItem = fileitem
	// 返回使用
	return ins, nil
}

// 关闭
func (ins *QueryInstance) Destroy() {
	// 释放文件控制
	if !ins.db.config.MemoryStorage {
		ins.db.releaseControlOfFile(ins)
	}
	// 清空数据
	ins.db = nil
	ins.key = nil
	ins.hash = nil
	ins.fileKey = ""
	ins.filePath = ""
	ins.searchHash = nil
	ins.targetFilePackage = nil
	ins.targetFileItem = nil
	ins.searchResultCache = nil
}
