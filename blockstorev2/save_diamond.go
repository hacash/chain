package blockstorev2

import (
	"encoding/binary"
	"github.com/hacash/core/stores"
)

// block data store
func (cs *BlockStore) SaveDiamond(diamond *stores.DiamondSmelt) error {
	// save
	query1, e1 := cs.diamondDB.CreateNewQueryInstance(diamond.Diamond)
	if e1 != nil {
		return e1
	}
	defer query1.Destroy()
	diamond_datas, e2 := diamond.Serialize()
	if e2 != nil {
		return e2
	}
	e3 := query1.Save(diamond_datas)
	if e3 != nil {
		return e3
	}
	// save diamond name by number key
	numberkey := make([]byte, 4)
	binary.BigEndian.PutUint32(numberkey, uint32(diamond.Number))
	query2, e4 := cs.diamondnumDB.CreateNewQueryInstance(numberkey)
	if e4 != nil {
		return e4
	}
	defer query2.Destroy()
	e5 := query2.Save(diamond.Diamond)
	if e5 != nil {
		return e5
	}
	return nil
}
