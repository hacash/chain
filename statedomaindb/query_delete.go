package statedomaindb

/**
 * clear search index cache
 */
func (ins *QueryInstance) Delete() error {

	// 内存数据库
	if ins.db.config.MemoryStorage {
		ins.db.MemoryStorageDB.Delete(ins.inputkey)
		return nil
	}
	// 磁盘数据库
	if ins.db.config.LevelDB {
		ins.db.LevelDB.Delete(ins.domainkey, nil)
		return nil
	}

	panic("NewStateDomainDB  must use LevelDB!")
}
