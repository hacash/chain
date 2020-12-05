package tinykvdb

import "encoding/binary"

func (kv *TinyKVDB) Set(key []byte, value []byte) error {

	if kv.UseLevelDB {
		return kv.ldb.Put(key, value, nil)
	}

	kv.wlock.Lock()
	defer kv.wlock.Unlock()

	hashkey := convertKeyToLen16Hash(key)
	// query
	query, e1 := kv.bashhashtreedb.CreateNewQueryInstance(hashkey)
	if e1 != nil {
		return e1
	}
	defer query.Destroy()
	// check val size
	if len(value) <= 8 {
		e2 := query.Save(append([]byte{0}, value...))
		if e2 != nil {
			return e2
		}
		return nil // save ok
	}
	// save to store file
	item, _ := query.Find()
	if item != nil && (item[0] == 1 || item[0] == ItemDelMark) {
		fstart := binary.BigEndian.Uint32(item[1:5])
		fvlen := binary.BigEndian.Uint32(item[5:9])
		if len(value) <= int(fvlen) {
			if len(value) < int(fvlen) {
				binary.BigEndian.PutUint32(item[5:9], uint32(len(value)))
				e2 := query.Save(item) // update val size
				if e2 != nil {
					return e2
				}
			}
			_, e2 := kv.storefile.WriteAt(value, int64(fstart))
			if e2 != nil {
				return e2
			}
		}
	}
	// append store
	stat, e3 := kv.storefile.Stat()
	if e3 != nil {
		return e3
	}
	file_size := stat.Size()
	_, e4 := kv.storefile.WriteAt(value, file_size)
	if e4 != nil {
		return e4
	}
	itemv := make([]byte, 9)
	itemv[0] = 1
	binary.BigEndian.PutUint32(itemv[1:5], uint32(file_size))
	binary.BigEndian.PutUint32(itemv[5:9], uint32(len(value)))
	e5 := query.Save(itemv)
	if e5 != nil {
		return e5
	}
	return nil
}
