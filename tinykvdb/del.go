package tinykvdb

func (kv *TinyKVDB) Del(key []byte) {
	hashkey := convertKeyToLen16Hash(key)
	// query
	query, _ := kv.bashhashtreedb.CreateNewQueryInstance(hashkey)
	if query == nil {
		return
	}
	defer query.Destroy()
	item, _ := query.Find()
	if item != nil {
		item[0] = ItemDelMark // delete mark
		query.Save(item)
		return
	}
}
