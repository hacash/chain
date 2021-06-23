package tinykvdb

import (
	"encoding/binary"
	"fmt"
)

func (kv *TinyKVDB) Get(key []byte) ([]byte, error) {

	if kv.UseLevelDB {
		v, _ := kv.GetOrCreateLevelDBwithPanic().Get(key, nil)
		if v != nil {
			return v, nil
		}
		return nil, nil
	}

	hashkey := convertKeyToLen16Hash(key)
	// query
	query, e1 := kv.bashhashtreedb.CreateNewQueryInstance(hashkey)
	if e1 != nil {
		return nil, e1
	}
	defer query.Destroy()
	value, e2 := query.Find()
	if e2 != nil {
		return nil, e2
	}
	if value == nil {
		return nil, nil // not find
	}
	if len(value) < 9 {
		return nil, fmt.Errorf("store file break down.")
	}
	if value[0] == 0 {
		return value[1:], nil // on tree
	}
	if value[0] == ItemDelMark {
		// delete mark
		return nil, nil // not find
	}
	valstart := binary.BigEndian.Uint32(value[1:5])
	vallength := binary.BigEndian.Uint32(value[5:9])
	// read from store file
	stat, e3 := kv.storefile.Stat()
	if e3 != nil {
		return nil, e3
	}
	// check size
	if stat.Size() < int64(valstart+vallength) {
		return nil, fmt.Errorf("store file size is error.")
	}
	resvalue := make([]byte, vallength)
	n, e4 := kv.storefile.ReadAt(resvalue, int64(valstart))
	if e4 != nil {
		return nil, e4
	}
	if n != int(vallength) {
		return nil, fmt.Errorf("ReadAt error")
	}
	// read successfully
	return resvalue, nil
}
