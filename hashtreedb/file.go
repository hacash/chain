package hashtreedb

/*
type TargetFilePackage struct {
	fileKey  string // Currently in use file key
	filePath string // Full path to the file currently in use without suffix

	gcFile    *os.File // . GC garbage collection file
	indexFile *os.File // . IDX index file
	dataFile  *os.File // . Dat data saving file

	//dataSegmentSize uint32 // 数据段尺寸
	//isStoreKeySize  uint32 // 如果储存了key，则key的尺寸   未储存则为0

}
*/

/*
// close and not use any more
func (tf *TargetFilePackage) Destroy() {
	tf.fileKey = ""
	tf.filePath = ""
	// Close file
	if tf.gcFile != nil {
		tf.gcFile.Close()
		tf.gcFile = nil
	}
	if tf.indexFile != nil {
		tf.indexFile.Close()
		tf.indexFile = nil
	}
	if tf.dataFile != nil {
		tf.dataFile.Close()
		tf.dataFile = nil
	}
}
*/

/*
// Waiting for file control (lock)
func (db *HashTreeDB) waitForTakeControlOfFile(ins *QueryInstance) (*lockFilePkgItem, error) {
	if len(ins.fileKey) == 0 {
		panic("(db *HashTreeDB) waitForTakeControlOfFile -> len(ins.fileKey) == 0")
	}

	// Operation lock
	db.filesOptLock.Lock()

	// Start detecting files
	var fileitem *lockFilePkgItem = nil
	fwlockptr := &sync.Mutex{}
	tarlock, ldok := db.filesWriteLock.Load(ins.fileKey)
	if ldok {
		fileitem = tarlock.(*lockFilePkgItem) // just get out
		fwlockptr = fileitem.lock
	}
	// Destination package
	var targetfilepkg = &TargetFilePackage{}
	targetfilepkg.fileKey = ins.fileKey
	targetfilepkg.filePath = ins.filePath
	var datfn = ins.filePath + ".dat"
	// Determine whether the file exists
	if fileitem != nil && fileitem.count >= 1 {
		// Cache is used and the file is not closed
		targetfilepkg = fileitem.targetFilePackageCache
		fileitem.count += 1 // Statistical quantity
	} else {
		// Cache does not exist
		// Check whether the file exists
		exist, e1 := PathExists(datfn)
		if e1 != nil {
			return nil, e1
		}
		if !exist {
			// Create directory if file does not exist
			basedir := path.Dir(datfn)
			e := os.MkdirAll(basedir, os.ModePerm)
			if e != nil {
				return nil, e
			}
		}
		// Create or open file
		err := openCreateTargetFiles(ins.filePath, targetfilepkg)
		if err != nil {
			return nil, err
		}
		// Save cache
		fileitem = &lockFilePkgItem{
			count:                  1,
			lock:                   fwlockptr,
			targetFilePackageCache: targetfilepkg,
		}
		db.filesWriteLock.Store(ins.fileKey, fileitem)
	}
	//db.existsFileKeys.Store(ins.fileKey, true) // 确定存在
	// Give package
	ins.targetFilePackage = targetfilepkg

	// Operation unlocking
	db.filesOptLock.Unlock()
	// Target file locking
	fwlockptr.Lock()

	// fwlockptr.Unlock()
	return fileitem, nil
}

func openCreateTargetFiles(fpfn string, targetfilepkg *TargetFilePackage) error {
	f1, e1 := os.OpenFile(fpfn+".dat", os.O_RDWR|os.O_CREATE, 0777) // |os.O_TRUNC =清空
	if e1 != nil {
		return e1
	}
	f2, e2 := os.OpenFile(fpfn+".idx", os.O_RDWR|os.O_CREATE, 0777)
	if e2 != nil {
		return e2
	}
	f2stat, ee2 := f2.Stat()
	if ee2 != nil {
		return ee2
	}
	if f2stat.Size() == 0 {
		// first
		f2.Write(bytes.Repeat([]byte{0}, IndexMenuSize))
	}
	f3, e3 := os.OpenFile(fpfn+".gc", os.O_RDWR|os.O_CREATE, 0777)
	if e3 != nil {
		return e3
	}
	targetfilepkg.dataFile = f1
	targetfilepkg.indexFile = f2
	targetfilepkg.gcFile = f3
	return nil
}

// Waiting for file control
func (db *HashTreeDB) releaseControlOfFile(ins *QueryInstance) error {
	ins.targetFileItem.count -= 1
	if ins.targetFileItem.count == 0 {
		// No one has adopted it. Close all files
		ins.targetFileItem.targetFilePackageCache.Destroy()
	}
	ins.targetFileItem.lock.Unlock() // Release file lock

	//fk := ins.fileKey

	return nil
}

*/

///////////////////////////////////////////////////////////////////////////////////
