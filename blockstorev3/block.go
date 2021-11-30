package blockstorev3

import (
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfacev3"
)

func (bs *BlockStore) SaveBlock(fullblock interfacev3.Block) error {

	ldb, e := bs.getDB()
	if e != nil {
		return e
	}

	blockhash := fullblock.Hash()
	blockbytes, e := fullblock.Serialize()
	if e != nil {
		return e
	}

	// save
	key := keyfix(blockhash, "block")
	return ldb.Put(key, blockbytes, nil)
}

func (bs *BlockStore) ReadBlockBytesByHash(blockhash fields.Hash) ([]byte, error) {

	ldb, e := bs.getDB()
	if e != nil {
		return nil, e
	}

	key := keyfix(blockhash, "block")
	blockbytes, e := ldb.Get(key, nil)
	if e != nil {
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
