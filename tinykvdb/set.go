package tinykvdb

import "encoding/binary"

func (kv *TinyKVDB) Set(key []byte, value []byte) {

	kv.wlock.Lock()
	defer kv.wlock.Unlock()

	hashkey := convertKeyToLen16Hash(key)
	// query
	query, _ := kv.bashhashtreedb.CreateNewQueryInstance(hashkey)
	if query == nil {
		return
	}
	defer query.Destroy()
	// check val size
	if len(value) <= 8 {
		query.Save(append([]byte{0}, value...))
		return
	}
	// save to store file
	item, _ := query.Find()
	if item != nil && (item[0] == 1 || item[0] == ItemDelMark) {
		fstart := binary.BigEndian.Uint32(item[1:5])
		fvlen := binary.BigEndian.Uint32(item[5:9])
		if len(value) <= int(fvlen) {
			if len(value) < int(fvlen) {
				binary.BigEndian.PutUint32(item[5:9], uint32(len(value)))
				query.Save(item) // update val size
			}
			kv.storefile.WriteAt(value, int64(fstart))
			return
		}
	}
	// append store
	stat, _ := kv.storefile.Stat()
	if stat == nil {
		return
	}
	file_size := stat.Size()
	kv.storefile.WriteAt(value, file_size)
	itemv := make([]byte, 9)
	itemv[0] = 1
	binary.BigEndian.PutUint32(itemv[1:5], uint32(file_size))
	binary.BigEndian.PutUint32(itemv[5:9], uint32(len(value)))
	query.Save(itemv)
}
