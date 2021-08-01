package blockstorev2

import (
	"encoding/binary"
	"github.com/hacash/core/fields"
)

// block data store
func (cs *BlockStore) ReadBlockBytesByHash(blkhash fields.Hash, readlen uint32) ([]byte, error) {
	blkdata, e1 := cs.blockdataDB.Read(blkhash, readlen)
	if e1 != nil {
		return nil, e1
	}
	return blkdata, nil
}

// block data store
// return: hash body error
func (cs *BlockStore) ReadBlockBytesByHeight(height uint64, readlen uint32) ([]byte, []byte, error) {
	numhash := make([]byte, 8)
	binary.BigEndian.PutUint64(numhash, height)
	// read
	query, e1 := cs.blknumhashDB.CreateNewQueryInstance(numhash)
	if e1 != nil {
		return nil, nil, e1
	}
	defer query.Destroy()
	blkhash, e2 := query.Find()
	if e2 != nil {
		return nil, nil, e2
	}
	if blkhash == nil || len(blkhash) != 32 {
		return nil, nil, nil // not find
	}
	// read
	resdata, err := cs.ReadBlockBytesByHash(blkhash, readlen)
	return blkhash, resdata, err
}

// block data store
func (cs *BlockStore) ReadBlockHeadBytesByHash(blkhash fields.Hash) ([]byte, error) {
	blkdata, _, e1 := cs.blockdataDB.ReadHead(blkhash)
	if e1 != nil {
		return nil, e1
	}
	return blkdata, nil
}

// block data store
func (cs *BlockStore) ReadBlockHeadBytesByHeight(height uint64) ([]byte, error) {
	// read
	blkhash, err := cs.ReadBlockHashByHeight(height)
	if err != nil {
		return nil, err
	}
	return cs.ReadBlockHeadBytesByHash(blkhash)
}

// block data store
func (cs *BlockStore) ReadBlockHashByHeight(height uint64) (fields.Hash, error) {
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
	return blkhash, nil
}
