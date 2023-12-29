package blockstorev3

import (
	"encoding/binary"
	"github.com/hacash/chain/leveldb"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
)

func (bs *BlockStore) ReadLastBlockHeight() (uint64, error) {
	if bs.lastBlockHeight > 0 {
		return bs.lastBlockHeight, nil
	}
	ldb, e := bs.getDB()
	if e != nil {
		return 0, e
	}
	vnum, e := ldb.Get([]byte("lastblkhei"), nil)
	if e != nil {
		return 0, e
	}
	return binary.BigEndian.Uint64(vnum), nil
}

func (bs *BlockStore) SaveBlock(fullblock interfaces.Block) error {

	ldb, e := bs.getDB()
	if e != nil {
		return e
	}

	blkhei := fullblock.GetHeight()
	blockhash := fullblock.Hash()
	blockbytes, e := fullblock.Serialize()
	if e != nil {
		return e
	}

	// save
	key := keyfix(blockhash, "block")
	e = ldb.Put(key, blockbytes, nil)
	if e != nil {
		return e
	}

	// save last height
	if blkhei > bs.lastBlockHeight {
		lstblkhei := make([]byte, 8)
		binary.BigEndian.PutUint64(lstblkhei, blkhei)
		e = ldb.Put([]byte("lastblkhei"), lstblkhei, nil)
		if e != nil {
			return e
		}
		bs.lastBlockHeight = blkhei
	}

	return nil
}

func (bs *BlockStore) ReadBlockBytesByHash(blockhash fields.Hash) ([]byte, error) {

	ldb, e := bs.getDB()
	if e != nil {
		return nil, e
	}

	key := keyfix(blockhash, "block")
	blockbytes, e := ldb.Get(key, nil)
	if e != nil {
		if e == leveldb.ErrNotFound {
			return nil, nil // not find
		}
		return nil, e
	}

	return blockbytes, nil
}

func (bs *BlockStore) ReadBlockBytesByHeight(blockheight uint64) (fields.Hash, []byte, error) {
	blkhx, e := bs.ReadBlockHashByHeight(blockheight)
	if e != nil {
		return nil, nil, e
	}

	blkbts, e := bs.ReadBlockBytesByHash(blkhx)
	if e != nil {
		return nil, nil, e
	}

	return blkhx, blkbts, nil
}

func (bs *BlockStore) ReadBlockHashByHeight(blockheight uint64) (fields.Hash, error) {

	ldb, e := bs.getDB()
	if e != nil {
		return nil, e
	}

	height := fields.BlockHeight(blockheight)
	heibts, e := height.Serialize()
	if e != nil {
		return nil, e
	}

	key := keyfix(heibts, "blkheihx")
	blockhash, e := ldb.Get(key, nil)
	if e != nil {
		if e == leveldb.ErrNotFound {
			return nil, nil // not find
		}
		return nil, e
	}

	return blockhash, nil
}

func (bs *BlockStore) UpdateSetBlockHashReferToHeight(blockheight uint64, blockhash fields.Hash) error {

	ldb, e := bs.getDB()
	if e != nil {
		return e
	}

	height := fields.BlockHeight(blockheight)
	heibts, e := height.Serialize()
	if e != nil {
		return e
	}

	key := keyfix(heibts, "blkheihx")
	return ldb.Put(key, blockhash, nil)
}
