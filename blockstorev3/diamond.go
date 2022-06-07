package blockstorev3

import (
	"github.com/hacash/chain/leveldb"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/stores"
)

func (bs *BlockStore) SaveDiamond(diaobj *stores.DiamondSmelt) error {
	dianame := []byte(diaobj.Diamond)
	diabts, e := diaobj.Serialize()
	if e != nil {
		return e
	}

	ldb, e := bs.getDB()
	if e != nil {
		return e
	}

	key := keyfix(dianame, "diamond")
	return ldb.Put(key, diabts, nil)
}

func (bs *BlockStore) ReadDiamond(name fields.DiamondName) (*stores.DiamondSmelt, error) {

	ldb, e := bs.getDB()
	if e != nil {
		return nil, e
	}

	key := keyfix(name, "diamond")
	databytes, e := ldb.Get(key, nil)
	if e != nil {
		if e == leveldb.ErrNotFound {
			return nil, nil // notfind
		}
		return nil, e
	}

	var diamondObj stores.DiamondSmelt
	_, e = diamondObj.Parse(databytes, 0)
	if e != nil {
		return nil, e
	}

	return &diamondObj, nil
}

func (bs *BlockStore) ReadDiamondNameByNumber(number uint32) (fields.DiamondName, error) {

	ldb, e := bs.getDB()
	if e != nil {
		return nil, e
	}

	height := fields.DiamondNumber(number)
	heibts, e := height.Serialize()
	if e != nil {
		return nil, e
	}

	key := keyfix(heibts, "dianum")
	dianame, e := ldb.Get(key, nil)
	if e != nil {
		if e == leveldb.ErrNotFound {
			return nil, nil // notfind
		}
		return nil, e
	}

	return dianame, nil
}

func (bs *BlockStore) ReadDiamondByNumber(number uint32) (*stores.DiamondSmelt, error) {

	dianame, e := bs.ReadDiamondNameByNumber(number)
	if e != nil {
		return nil, e
	}

	return bs.ReadDiamond(dianame)
}

// Set the diamond name pointed by the diamond number
func (bs *BlockStore) UpdateSetDiamondNameReferToNumber(number uint32, dianame fields.DiamondName) error {

	ldb, e := bs.getDB()
	if e != nil {
		return e
	}

	height := fields.DiamondNumber(number)
	heibts, e := height.Serialize()
	if e != nil {
		return e
	}

	key := keyfix(heibts, "dianum")
	return ldb.Put(key, dianame, nil)
}
