package hashtreedb

import (
	"bytes"
	"os"
	"path"
	"strings"
	"sync"
)



type TargetFilePackage struct {

	fileKey  string  // 当前正在使用的文件 key
	filePath string  // 当前正在使用的文件 完整路径但不带后缀名

	gcFile       *os.File // .gc  垃圾回收文件
	indexFile    *os.File // .idx 索引文件
	dataFile     *os.File // .dat 数据保存文件

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
func (db *HashTreeDB) waitForTakeControlOfFile(ins *QueryInstance) (*sync.Mutex, error) {
	if len(ins.fileKey) == 0 {
		panic("(db *HashTreeDB) waitForTakeControlOfFile -> len(ins.fileKey) == 0")
	}
	db.fileOptLock.Lock()
	fwlockptr := &sync.Mutex{}
	tarlock, ldok := db.fileWriteLockMutexs.LoadOrStore(ins.fileKey, fwlockptr)
	if ldok {
		fwlockptr = tarlock.(*sync.Mutex) // just get out
	}
	db.fileOptLock.Unlock()
	fwlockptr.Lock() // 目标文件加锁
	// 目标文件包
	var targetfilepkg = &TargetFilePackage{}
	targetfilepkg.fileKey = ins.fileKey
	targetfilepkg.filePath = ins.filePath
	var datfn = ins.filePath + ".dat"
	// 判断文件是否存在
	_, is_exist := db.existsFileKeys.Load(ins.fileKey)
	if ! is_exist {
		//fmt.Println("-- PathExists(datfn) --")
		// 从文件夹检查文件是否存在
		exist, e1 := PathExists(datfn)
		if e1 != nil {
			return nil, e1
		}
		is_exist = exist
	}else{
		//fmt.Println("is_exist  ==  true")
	}
	// 创建目标夹和文件
	if ! is_exist {
		basedir := path.Dir(datfn)
		os.MkdirAll(basedir, os.ModePerm)
		err := openCreateTargetFiles(ins.filePath, targetfilepkg)
		if err != nil {
			return nil, err
		}
		//fmt.Println("+++++++++++++++++++++ create file")
	}else{
		// 检查缓存池
		if db.targetFilePackagePool != nil && strings.Compare(ins.fileKey, db.targetFilePackagePool.fileKey) == 0 {
			targetfilepkg = db.targetFilePackagePool // 直接使用缓存
			//fmt.Println("-------------------- targetFilePackagePool")
		}else{
			// 新打开文件并创建
			err := openCreateTargetFiles(ins.filePath, targetfilepkg)
			if err != nil {
				return nil, err
			}
			//fmt.Println("=================== open file")
		}
	}
	if db.targetFilePackagePool == nil {
		db.targetFilePackagePool = targetfilepkg // 暂保留一条缓存
	}
	db.existsFileKeys.Store(ins.fileKey,true) // 确定存在
	// 给出文件包
	ins.targetFilePackage = targetfilepkg
	// fwlockptr.Unlock()
	return fwlockptr, nil
}



func openCreateTargetFiles(fpfn string, targetfilepkg *TargetFilePackage) error {
	f1, e1 := os.OpenFile(fpfn + ".dat", os.O_RDWR|os.O_CREATE, 0777) // |os.O_TRUNC =清空
	if e1 != nil {
		return e1
	}
	f2, e2 := os.OpenFile(fpfn + ".idx", os.O_RDWR|os.O_CREATE, 0777)
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
	f3, e3 := os.OpenFile(fpfn + ".gc", os.O_RDWR|os.O_CREATE, 0777)
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
	ins.targetFileWriteJustUnlock.Unlock() // 释放文件锁

	//fk := ins.fileKey


	return nil
}








///////////////////////////////////////////////////////////////////////////////////


