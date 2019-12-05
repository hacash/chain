package tinykvdb

import "github.com/hacash/chain/hashtreedb"

/**
 * small kv db
 */
type TinyKVDB struct {

	bashhashtreedb *hashtreedb.HashTreeDB

}




// create DataBase
func NewTinyKVDB() *TinyKVDB {
	db := &TinyKVDB{
	}
	return db
}


















