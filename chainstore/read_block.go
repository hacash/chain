package chainstore

import (
	"encoding/binary"
	"github.com/hacash/core/fields"
)

// block data store
func (cs *ChainStore) ReadBlockBytesByHash(blkhash fields.Hash) ([]byte, error) {
	blkdata, e1 := cs.blockdataDB.Read(blkhash)
	if e1 != nil {
		return nil, e1
	}
	return blkdata, nil
}

// block data store
func (cs *ChainStore) ReadBlockBytesByHeight(height uint64) ([]byte, error) {
	numhash := make([]byte, 8)
	binary.BigEndian.PutUint64(numhash, height)
	// read
	query, e1 := cs.blknumhashDB.CreateNewQueryInstance(numhash)
	if e1 != nil {
		return nil, e1
	}
	defer query.Destroy()
	blkhash, e2 := query.Find()
	if e2 != nil {
		return nil, e2
	}
	// read
	return cs.ReadBlockBytesByHash(blkhash)
}

// block data store
func (cs *ChainStore) ReadBlockHeadBytesByHash(blkhash fields.Hash) ([]byte, error) {
	blkdata, _, e1 := cs.blockdataDB.ReadHead(blkhash)
	if e1 != nil {
		return nil, e1
	}
	return blkdata, nil
}

// block data store
func (cs *ChainStore) ReadBlockHeadBytesByHeight(height uint64) ([]byte, error) {
	numhash := make([]byte, 8)
	binary.BigEndian.PutUint64(numhash, height)
	// read
	query, e1 := cs.blknumhashDB.CreateNewQueryInstance(numhash)
	if e1 != nil {
		return nil, e1
	}
	defer query.Destroy()
	blkhash, e2 := query.Find()
	if e2 != nil {
		return nil, e2
	}
	// read
	return cs.ReadBlockHeadBytesByHash(blkhash)
}
