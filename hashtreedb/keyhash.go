package hashtreedb

import ()

/*
// find target file
func (db *HashTreeDB) locateTargetFilePath(hash []byte) (string, string, []byte) {
	use_hash, _, pathary := db.spreadHashToIndexPath(hash)
	mpath := strings.Join(pathary, "/")
	filepath := strings.TrimRight(db.config.FileAbsPath, "/") + "/" + mpath + db.config.FileName // + ".dat"
	return filepath, mpath + "-", use_hash
}

// Split hash as path
func (db *HashTreeDB) spreadHashToIndexPath(hash []byte) ([]byte, []byte, []string) {
	lv := int(db.config.FileDividePartitionLevel)
	if lv > 0 {
		path := make([]string, lv, lv)
		for i := 0; i < lv; i++ {
			path[i] = strconv.Itoa(int(hash[i]))
		}
		return hash[lv:], hash[:lv], path
	} else {
		return hash, []byte{}, []string{}
	}
}

// Convert the original key expansion to hash value
func (db *HashTreeDB) convertKeyToHash(key []byte) []byte {
	// reverse
	if db.config.KeyReverse {
		key = ReverseHashOrder(key) // copy
	} else {
		var hsdt = make([]byte, len(key))
		copy(hsdt, key)
		key = hsdt // copy
	}
	// key prefix supplement
	if db.config.KeyPrefixSupplement > 0 && db.config.KeyPrefixSupplement < 30 {
		sha256hash := sha256.New()
		sha256hash.Write(key)
		prefixhash := sha256hash.Sum(nil)
		buf := bytes.NewBuffer(prefixhash[0:db.config.KeyPrefixSupplement])
		buf.Write(key)
		key = buf.Bytes()
	}
	// ok
	hexkey := hex.EncodeToString(key)
	//fmt.Println(hexkey)
	realusekey := make([]byte, len(hexkey))
	for i := 0; i < len(hexkey); i++ {
		ha := hexkey[i]
		if ha > 57 {
			realusekey[i] = ha - 97 + 10
		} else {
			realusekey[i] = ha - 48
		}
	}
	//fmt.Println(realusekey)
	return realusekey
}

*/
