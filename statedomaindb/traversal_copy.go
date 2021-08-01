package statedomaindb

// 遍历拷贝、修改、删除数据， 只能是单文件数据库

func (this *StateDomainDB) TraversalCopy(target *StateDomainDB) error {
	// 不能把文件数据库的内容，拷贝到内存数据库

	// 内存数据库
	if target.config.MemoryStorage {
		// 遍历
		target.MemoryStorageDB.wlok.Lock()
		defer target.MemoryStorageDB.wlok.Unlock()
		for k, v := range target.MemoryStorageDB.Datas {
			//fmt.Println("TraversalCopy", fields.Address([]byte(k)).ToReadable(), v)
			distins, e0 := newQueryInstance(this, []byte(k))
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

	panic("NewStateDomainDB TraversalCopy cannot must use LevelDB!")

}
