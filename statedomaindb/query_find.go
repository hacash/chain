package statedomaindb

import ()

/**
 * search index exist
 */
func (ins *QueryInstance) Exist() (bool, error) {
	// In memory database
	if ins.db.config.MemoryStorage {
		return ins.db.MemoryStorageDB.Exist(ins.inputkey), nil
	}

	// LevelDB
	if ins.db.config.LevelDB {
		val, err := ins.db.LevelDB.Get(ins.domainkey, nil)
		if err != nil || val == nil {
			return false, nil // error
		}
		return true, nil
	}

	panic("NewStateDomainDB  must use LevelDB!")

}

/**
 * search index file and get the item part
 */
func (ins *QueryInstance) Find() ([]byte, error) {
	// In memory database
	if ins.db.config.MemoryStorage {
		val, ok := ins.db.MemoryStorageDB.Read(ins.inputkey)
		if !ok || val == nil {
			return nil, nil
		}
		// copy
		if ins.db.config.SupplementalMaxValueSize > 0 {
			retdts := make([]byte, ins.db.config.SupplementalMaxValueSize) // 补充不足的长度
			copy(retdts, val)
			//fmt.Println("MemoryStorageDB Find", fields.Address(ins.domainkey).ToReadable(), retdts)
			return retdts, nil
		}
		// Original stored data
		return val, nil
	}

	// LevelDB
	if ins.db.config.LevelDB {
		val, err := ins.db.LevelDB.Get(ins.domainkey, nil)
		if err != nil || val == nil {
			return nil, nil // error or not find
		}
		// copy
		if ins.db.config.SupplementalMaxValueSize > 0 {
			retdts := make([]byte, ins.db.config.SupplementalMaxValueSize) // 补充不足的长度
			copy(retdts, val)
			//fmt.Println("LevelDB Find", fields.Address(ins.domainkey).ToReadable(), retdts)
			return retdts, nil
		}
		// Original stored data
		return val, nil
	}

	panic("NewStateDomainDB must use LevelDB!")
}
