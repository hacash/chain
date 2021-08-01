package statedomaindb

// Get
func (db *StateDomainDB) Get(key []byte) ([]byte, error) {

	ins, e3 := db.CreateNewQueryInstance(key)
	if e3 != nil {
		return nil, e3
	}
	vdatas, e4 := ins.Find()
	if e4 != nil {
		return nil, e4
	}
	// ok
	return vdatas, nil
}

// Set
func (db *StateDomainDB) Set(key []byte, val []byte) error {
	ins, e3 := db.CreateNewQueryInstance(key)
	if e3 != nil {
		return e3
	}
	e4 := ins.Save(val)
	if e4 != nil {
		return e4
	}
	// ok
	return nil
}
