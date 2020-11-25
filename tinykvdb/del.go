package tinykvdb

func (kv *TinyKVDB) Del(key []byte) error {
	hashkey := convertKeyToLen16Hash(key)
	// query
	query, e1 := kv.bashhashtreedb.CreateNewQueryInstance(hashkey)
	if e1 != nil {
		return e1
	}
	defer query.Destroy()
	item, e2 := query.Find()
	if e2 != nil {
		return e2
	}
	if item != nil {
		item[0] = ItemDelMark // delete mark
		e1 := query.Save(item)
		if e1 != nil {
			return e1
		}
	}
	return nil // ok
}
