package statedomaindb

import "bytes"

// Query instance

type QueryInstance struct {
	db *StateDomainDB

	inputkey  []byte
	domainkey []byte
}

func newQueryInstance(db *StateDomainDB, inkey []byte) (*QueryInstance, error) {
	// Suffix key
	keybuf := bytes.NewBuffer(inkey)
	keybuf.Write([]byte(db.config.KeyDomainName))
	// establish
	ins, e := newQueryInstanceByRealUseKey(db, keybuf.Bytes())
	if e != nil {
		return nil, e
	}
	ins.inputkey = inkey
	return ins, nil
}

func newQueryInstanceByRealUseKey(db *StateDomainDB, realkey []byte) (*QueryInstance, error) {
	ins := &QueryInstance{
		db:        db,
		domainkey: realkey,
	}
	// If it is an in memory database, do not open the local file
	if db.config.MemoryStorage {
		return ins, nil
	}
	// If it is level dB, do not open the file
	if db.config.LevelDB {
		return ins, nil
	}

	panic("must use leveldb")

}

// close
func (ins *QueryInstance) Destroy() {
	// Release file control
	if !ins.db.config.MemoryStorage && !ins.db.config.LevelDB {
		// ins.db.releaseControlOfFile(ins)
	}
	// wipe data 
	ins.db = nil
	ins.inputkey = nil
	ins.domainkey = nil
}
