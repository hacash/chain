package statedomaindb

import (
	"fmt"
)

//func (ins *QueryInstance) Save(valuedatas []byte) (ValueSegmentOffset uint32, err error) {
func (ins *QueryInstance) Save(valuedatas []byte) error {

	// In memory database
	if ins.db.config.MemoryStorage {
		// copy
		retdts := make([]byte, len(valuedatas))
		copy(retdts, valuedatas)
		//fmt.Println("MemoryStorageDB Save", fields.Address(ins.domainkey).ToReadable(), valuedatas)
		ins.db.MemoryStorageDB.Save(ins.inputkey, retdts)
		return nil
	}

	// LevelDB
	if ins.db.config.LevelDB {
		err := ins.db.LevelDB.Put(ins.domainkey, valuedatas, nil)
		if err != nil {
			return fmt.Errorf("QueryInstance.LevelDB.Save Error: %s", err.Error())
		}
		return nil
	}

	panic("NewStateDomainDB  must use LevelDB!")

	return nil
}
