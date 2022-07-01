package hashtreedb

// Query instance

type QueryInstance struct {
	db *HashTreeDB

	key []byte
	//hash       []byte
	//fileKey    string
	//filePath   string
	//searchHash []byte

	//targetFilePackage *TargetFilePackage
	//targetFileItem    *lockFilePkgItem

	// search cache
	//searchResultCache *FindValueOffsetItem
}

/**
 * QueryInstance {
 *     Find()
 *     Save(valuebytes []byte)
 * }
 */
func newQueryInstance(db *HashTreeDB, key []byte) (*QueryInstance, error) {

	ins := &QueryInstance{
		db:  db,
		key: key,
		//searchResultCache: nil,
	}
	// If it is an in memory database, do not open the local file
	if db.config.MemoryStorage {
		return ins, nil
	}
	// If it is level dB, do not open the file
	if db.config.LevelDB {
		return ins, nil
	}

	panic("must use leveldb")

	/*
		ins.hash = db.convertKeyToHash(key)
		ins.filePath, ins.fileKey, ins.searchHash = db.locateTargetFilePath(ins.hash)
		//fmt.Println("newQueryInstance searchHash ", ins.searchHash)
		// Waiting for file control
		fileitem, err := db.waitForTakeControlOfFile(ins)
		if err != nil {
			return nil, err
		}
		ins.targetFileItem = fileitem
		// Return to use
	*/
	return ins, nil
}

// close
func (ins *QueryInstance) Destroy() {
	// Release file control
	if !ins.db.config.MemoryStorage && !ins.db.config.LevelDB {
		// ins.db.releaseControlOfFile(ins)
	}
	// wipe data
	ins.db = nil
	ins.key = nil
	//ins.hash = nil
	//ins.fileKey = ""
	//ins.filePath = ""
	//ins.searchHash = nil
	//ins.targetFilePackage = nil
	//ins.targetFileItem = nil
	//ins.searchResultCache = nil
}
