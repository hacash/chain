package hashtreedb

import (
	"bytes"
	"os"
	"path"
	"sync"
)

type TargetFilePackage struct {
	fileKey  string // 当前正在使用的文件 key
	filePath string // 当前正在使用的文件 完整路径但不带后缀名

	gcFile    *os.File // .gc  垃圾回收文件
	indexFile *os.File // .idx 索引文件
	dataFile  *os.File // .dat 数据保存文件

	//dataSegmentSize uint32 // 数据段尺寸
	//isStoreKeySize  uint32 // 如果储存了key，则key的尺寸   未储存则为0

}

// close and not use any more
func (tf *TargetFilePackage) Destroy() {
	tf.fileKey = ""
	tf.filePath = ""
	// 关闭文件
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

// 等待获取文件控制权（锁）
func (db *HashTreeDB) waitForTakeControlOfFile(ins *QueryInstance) (*lockFilePkgItem, error) {
	if len(ins.fileKey) == 0 {
		panic("(db *HashTreeDB) waitForTakeControlOfFile -> len(ins.fileKey) == 0")
	}

	// 操作锁定
	db.filesOptLock.Lock()

	// 开始检测文件
	var fileitem *lockFilePkgItem = nil
	fwlockptr := &sync.Mutex{}
	tarlock, ldok := db.filesWriteLock.Load(ins.fileKey)
	if ldok {
		fileitem = tarlock.(*lockFilePkgItem) // just get out
		fwlockptr = fileitem.lock
	}
	// 目标文件包
	var targetfilepkg = &TargetFilePackage{}
	targetfilepkg.fileKey = ins.fileKey
	targetfilepkg.filePath = ins.filePath
	var datfn = ins.filePath + ".dat"
	// 判断文件是否存在
	if fileitem != nil && fileitem.count >= 1 {
		// 使用缓存，并且文件没有被关闭
		targetfilepkg = fileitem.targetFilePackageCache
		fileitem.count += 1 // 统计数量
	} else {
		// 缓存不存在
		// 检查文件是否存在
		exist, e1 := PathExists(datfn)
		if e1 != nil {
			return nil, e1
		}
		if !exist {
			// 文件不存在，则创建目录
			basedir := path.Dir(datfn)
			e := os.MkdirAll(basedir, os.ModePerm)
			if e != nil {
				return nil, e
			}
		}
		// 创建或打开文件
		err := openCreateTargetFiles(ins.filePath, targetfilepkg)
		if err != nil {
			return nil, err
		}
		// 保存缓存
		fileitem = &lockFilePkgItem{
			count:                  1,
			lock:                   fwlockptr,
			targetFilePackageCache: targetfilepkg,
		}
		db.filesWriteLock.Store(ins.fileKey, fileitem)
	}
	//db.existsFileKeys.Store(ins.fileKey, true) // 确定存在
	// 给出文件包
	ins.targetFilePackage = targetfilepkg

	// 操作解锁
	db.filesOptLock.Unlock()
	// 目标文件加锁
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

// 等待获取文件控制权
func (db *HashTreeDB) releaseControlOfFile(ins *QueryInstance) error {
	ins.targetFileItem.count -= 1
	if ins.targetFileItem.count == 0 {
		// 已经无人采用，关闭所有文件
		ins.targetFileItem.targetFilePackageCache.Destroy()
	}
	ins.targetFileItem.lock.Unlock() // 释放文件锁

	//fk := ins.fileKey

	return nil
}

///////////////////////////////////////////////////////////////////////////////////
