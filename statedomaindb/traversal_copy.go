package statedomaindb

// Traverse, copy, modify and delete data, only single file database

func (this *StateDomainDB) TraversalCopy(target *StateDomainDB) error {
	// The contents of the file database cannot be copied to the memory database

	// In memory database
	if target.config.MemoryStorage {
		// ergodic
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
