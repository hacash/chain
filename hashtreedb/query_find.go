package hashtreedb

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

/**
 * search index exist
 */
func (ins *QueryInstance) Exist() (bool, error) {
	ins.ClearSearchIndexCache()
	ofstItem, err := ins.SearchIndex()
	if err != nil {
		return false, err // error
	}
	if ofstItem == nil {
		return false, nil // not find
	}
	if bytes.Compare(ins.key, ofstItem.ValueKey) == 0 {
		// find ok
		return true, nil
	} else {
		// other one
		return false, nil
	}
}

/**
 * search index file and get the item part
 */
func (ins *QueryInstance) Find() ([]byte, error) {
	ins.ClearSearchIndexCache()
	ofstItem, err := ins.SearchIndex()
	if err != nil {
		return nil, err // error
	}
	if ofstItem == nil {
		return nil, nil // not find
	}
	if ofstItem.Type == 0 {
		return nil, nil // not find
	}
	e2 := ins.readSegmentDataFillItem(ofstItem, true)
	if e2 != nil {
		return nil, e2 // error
	}
	if bytes.Compare(ins.key, ofstItem.ValueKey) == 0 {
		// read target ok
		return ofstItem.ValueBody, nil
	} else {
		// other one
		return nil, nil
	}
}

/**
 * search index file and get the item part
 */
func (ins *QueryInstance) readSegmentDataFillItem(fditem *FindValueOffsetItem, isreadvalue bool) error {
	// read data from file
	readsz := int(ins.db.config.segmentValueSize)
	if !isreadvalue {
		readsz -= int(ins.db.config.MaxValueSize)
	}
	var databytes = make([]byte, readsz)
	var rdoffset = fditem.ValueSegmentOffset * ins.db.config.segmentValueSize
	rdsz, rderr := ins.targetFilePackage.dataFile.ReadAt(databytes, int64(rdoffset))
	if rderr != nil {
		return rderr // return error
	}
	if rdsz < readsz {
		return fmt.Errorf("index file breakdown.")
	}
	var mksz int = 0
	if ins.db.config.SaveMarkBeforeValue {
		mksz = 1
		fditem.ValueMark = databytes[0]
	}
	fditem.ValueKey = databytes[mksz : mksz+int(ins.db.config.KeySize)]
	if isreadvalue {
		fditem.ValueBody = databytes[mksz+int(ins.db.config.KeySize):]
	} else {
		fditem.ValueBody = nil
	}
	return nil
}

/**
 * clear search index cache
 */
func (ins *QueryInstance) ClearSearchIndexCache() {
	ins.searchResultCache = nil
}

/**
 * search index file and get the item part
 */
func (ins *QueryInstance) SearchIndex() (*FindValueOffsetItem, error) {
	//
	if len(ins.searchHash) == 0 {
		panic("QueryInstance searchHash is null.")
	}
	// check cache
	if ins.searchResultCache != nil {
		return ins.searchResultCache, nil // cache
	}
	// seek file
	var idxrofst uint32 = 0
	var qhash_si = 0 // already drop file part prefix from searchHash
	for {
		// 例外
		if qhash_si >= len(ins.searchHash) {
			return nil, fmt.Errorf("search to the final.")
		}
		var curidxseg = make([]byte, IndexMenuSize)
		rdsz, rderr := ins.targetFilePackage.indexFile.ReadAt(curidxseg, int64(idxrofst))
		if rderr != nil {
			//stat, _ := ins.targetFilePackage.indexFile.Stat()
			//fmt.Println(idxrofst, stat.Size())
			return nil, rderr // return error
		}
		if rdsz == 0 {
			return nil, nil // file empty not find
		}
		if rdsz != IndexMenuSize {
			return nil, fmt.Errorf("read index file length is not 'HashTreeDBMenuSize'.")
		}
		itofst := IndexItemSize * int(ins.searchHash[qhash_si])
		var index_item = curidxseg[itofst : itofst+IndexItemSize]
		var item_type = index_item[0]
		ffdok := &FindValueOffsetItem{
			Type:                       item_type,
			searchHash:                 ins.searchHash,
			searchCount:                qhash_si,
			IndexMenuSelfSegmentOffset: uint32(idxrofst) / uint32(IndexMenuSize),
			IndexItemSelfAlignment:     uint32(itofst),
		}
		if item_type == IndexItemTypeNull {
			ffdok.ValueSegmentOffset = 0
			ffdok.ValueKey = nil
			ffdok.ValueBody = nil
			return ffdok, nil // not find
		} else if item_type == IndexItemTypeValue || item_type == IndexItemTypeValueDelete {
			ffdok.ValueSegmentOffset = binary.BigEndian.Uint32(index_item[1:5])
			return ffdok, nil // !!!!!! SUCCESS FIND !!!!!!
		} else if item_type == IndexItemTypeBranch {
			// next step to check menu
			menusegnum := binary.BigEndian.Uint32(index_item[1:5])
			idxrofst = menusegnum * uint32(IndexMenuSize)
			qhash_si++
		} else {
			return nil, fmt.Errorf("index file breakdown.")
		}
	}
	return nil, nil
}
