package statedomaindb

/**
 * clear search index cache
 */
func (ins *QueryInstance) Delete() error {

	// In memory database
	if ins.db.config.MemoryStorage {
		ins.db.MemoryStorageDB.Delete(ins.inputkey)
		return nil
	}
	// Disk database
	if ins.db.config.LevelDB {
		ins.db.LevelDB.Delete(ins.domainkey, nil)
		return nil
	}

	panic("NewStateDomainDB  must use LevelDB!")
}
