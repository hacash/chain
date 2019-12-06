package tinykvdb

import "encoding/binary"

func (kv *TinyKVDB) Get(key []byte) []byte {
	hashkey := convertKeyToLen16Hash(key)
	// query
	query, _ := kv.bashhashtreedb.CreateNewQueryInstance(hashkey)
	if query == nil {
		return nil
	}
	defer query.Destroy()
	value, _ := query.Find()
	if len(value) < 9 {
		return nil
	}
	if value[0] == 0 {
		return value[1:] // on tree
	}
	if value[0] == ItemDelMark {
		// delete mark
		return nil
	}
	valstart := binary.BigEndian.Uint32(value[1:5])
	vallength := binary.BigEndian.Uint32(value[5:9])
	// read from store file
	stat, _ := kv.storefile.Stat()
	if stat == nil {
		return nil
	}
	// check size
	if stat.Size() < int64(valstart+vallength) {
		return nil
	}
	resvalue := make([]byte, vallength)
	n, _ := kv.storefile.ReadAt(resvalue, int64(valstart))
	if n != int(vallength) {
		return nil
	}
	// read successfully
	return resvalue
}
