package blockstorev2

import (
	"encoding/binary"
	"fmt"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/stores"
)

// block data store
func (cs *BlockStore) ReadDiamond(diamond_name fields.DiamondName) (*stores.DiamondSmelt, error) {
	// find
	query1, e1 := cs.diamondDB.CreateNewQueryInstance(diamond_name)
	if e1 != nil {
		return nil, e1
	}
	defer query1.Destroy()
	diamonddatas, e3 := query1.Find()
	if e3 != nil {
		return nil, e3
	}
	if len(diamonddatas) == 0 {
		return nil, fmt.Errorf("diamond store file break.")
	}
	var diamond stores.DiamondSmelt
	_, e4 := diamond.Parse(diamonddatas, 0)
	if e4 != nil {
		return nil, e4
	}
	// ok
	return &diamond, nil
}

// block data store
func (cs *BlockStore) ReadDiamondNameByNumber(number uint32) (fields.DiamondName, error) {
	// find by number key
	numberkey := make([]byte, 4)
	binary.BigEndian.PutUint32(numberkey, number)
	query1, e1 := cs.diamondnumDB.CreateNewQueryInstance(numberkey)
	if e1 != nil {
		return nil, e1
	}
	defer query1.Destroy()
	diamondnamedatas, e3 := query1.Find()
	if e3 != nil {
		return nil, e3
	}
	if len(diamondnamedatas) < 6 {
		return nil, fmt.Errorf("diamond num store file break.")
	}
	// find by name
	return diamondnamedatas, nil
}

// block data store
func (cs *BlockStore) ReadDiamondByNumber(number uint32) (*stores.DiamondSmelt, error) {
	// find by number key
	dmdname, e1 := cs.ReadDiamondNameByNumber(number)
	if e1 != nil {
		return nil, e1
	}
	// find by name
	return cs.ReadDiamond(dmdname)
}
