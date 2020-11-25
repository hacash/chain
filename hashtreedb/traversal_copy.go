package hashtreedb

import (
	"fmt"
	"os"
)

// 遍历拷贝、修改、删除数据， 只能是单文件数据库

func (this *HashTreeDB) TraversalCopy(target *HashTreeDB, datafilemustexist bool) error {
	// 不能把文件数据库的内容，拷贝到内存数据库

	// 内存数据库
	if target.config.MemoryStorage {
		// 遍历
		target.MemoryStorageDB.wlok.Lock()
		defer target.MemoryStorageDB.wlok.Unlock()
		for k, v := range target.MemoryStorageDB.Datas {
			//fmt.Println("TraversalCopy", fields.Address([]byte(k)).ToReadable(), v)
			distins, e0 := this.CreateNewQueryInstance([]byte(k))
			if e0 != nil {
				return e0
			}
			if v.IsDelete {
				e0 = distins.Delete()
				if e0 != nil {
					distins.Destroy()
					return e0
				}
			} else {
				e0 = distins.Save(v.Value)
				if e0 != nil {
					distins.Destroy()
					return e0
				}
			}
			distins.Destroy()
			// next one
		}
		return nil
	}

	// 文件数据库
	if target.config.FileDividePartitionLevel > 0 {
		return fmt.Errorf("unsupported operations for TraversalCopy: config.FilePartitionLevel must be 0")
	}
	if target.config.ForbidGC != true {
		return fmt.Errorf("unsupported operations for TraversalCopy: config.ForbidGC must be true")
	}
	if target.config.SaveMarkBeforeValue != true {
		return fmt.Errorf("unsupported operations for TraversalCopy: config.SaveMarkBeforeValue must be true")
	}
	filepath, _, _ := target.locateTargetFilePath([]byte{})
	datafilename := filepath + ".dat"
	datafile, fe := os.OpenFile(datafilename, os.O_RDWR|os.O_CREATE, 0777)
	if fe != nil {
		if datafilemustexist {
			return fmt.Errorf("unsupported operations for TraversalCopy: file '" + datafilename + "' must be existence")
		} else {
			return nil // not hav any data
		}
	}
	defer datafile.Close()
	datafilestat, se := datafile.Stat()
	if se != nil {
		return se
	}
	datafilesize := datafilestat.Size()
	if datafilesize == 0 {
		return nil // empty
	}
	if datafilesize%int64(target.config.segmentValueSize) != 0 {
		return fmt.Errorf("data file break down.")
	}
	// copy
	onereadsize := uint32(4096)
	onereadsize = onereadsize / target.config.segmentValueSize * target.config.segmentValueSize
	datafileseek := int64(0)
	for {
		if datafileseek >= datafilesize {
			return nil // end
		}
		datasegments := make([]byte, onereadsize)
		rdlen, re := datafile.ReadAt(datasegments, datafileseek)
		if rdlen == 0 && re != nil {
			return re
		}
		if rdlen == 0 {
			return nil // end
		}
		if rdlen%int(target.config.segmentValueSize) != 0 {
			return fmt.Errorf("index file break down.")
		}
		// do copy
		err := this.recursTraversalCopy(datasegments[0:rdlen], target)
		if err != nil {
			return err
		}
		// ok next
		datafileseek += int64(onereadsize)
	}

}

func (this *HashTreeDB) recursTraversalCopy(values_list []byte, target *HashTreeDB) error {
	segmentValueSize := int(target.config.segmentValueSize)
	for i := 0; i < len(values_list)/segmentValueSize; i++ {
		seek := i * segmentValueSize
		value := values_list[seek : seek+segmentValueSize]
		markty := value[0]
		if markty != IndexItemTypeValue && markty != IndexItemTypeValueDelete {
			continue
		}
		onekey := value[1 : 1+target.config.KeySize]
		// set
		query, e := this.CreateNewQueryInstance(onekey)
		if e != nil {
			return e
		}
		if markty == IndexItemTypeValue {
			onevaluebody := value[1+target.config.KeySize:]
			e = query.Save(onevaluebody) // set
			if e != nil {
				query.Destroy()
				return e
			}
		} else {
			e = query.Delete() // del
			if e != nil {
				query.Destroy()
				return e
			}
		}
		query.Destroy()
		// OK NEXT
	}
	return nil
}
