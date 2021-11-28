package blockstorev3

import (
	leveldb "github.com/hacash/chain/leveldb"
)

type BlockStore struct {

	// config
	config *BlockStoreConfig

	// level db
	ldb *leveldb.DB
}
