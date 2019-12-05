package biglogdb

import (
	"github.com/hacash/chain/hashtreedb"
	"path"
)


const (
	// filenum + fileseek + length + space
	BlockDataPtrSize = 4 + 4 + 4 + 4
)




/**
 * config
 */
type BigLogDBConfig struct {
	DataDir string
	KeySize uint8
	LogHeadMaxSize int
	BlockPartFileMaxSize int64
	FileDividePartitionLevel uint8
}



func NewBigLogDBConfig(
	DataDir string,
	keySize uint8,
) *BigLogDBConfig {
	return &BigLogDBConfig{
		DataDir: DataDir,
		KeySize: keySize,
		LogHeadMaxSize: 0,
		FileDividePartitionLevel: 0,
	}
}





/**
 * big log db
 */
type BigLogDB struct {

	config *BigLogDBConfig

	bashhashtreedb *hashtreedb.HashTreeDB


}




// create DataBase
func NewBigLogDB(config *BigLogDBConfig) *BigLogDB {
	hsdbdir := path.Join(config.DataDir, "INDEXS")
	hsdbcnf := hashtreedb.NewHashTreeDBConfig(
		hsdbdir,
		uint32(config.LogHeadMaxSize) + BlockDataPtrSize,
		config.KeySize,
		)
	// copy cnf con
	hsdbcnf.FileDividePartitionLevel = config.FileDividePartitionLevel
	// new tree db
	basedb := hashtreedb.NewHashTreeDB(hsdbcnf)
	db := &BigLogDB{
		config: config,
		bashhashtreedb: basedb,
	}
	return db
}











