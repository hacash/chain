package hashtreedb

import (
	"bytes"
	"fmt"
)

/**
 * append value to file
 */
func (ins *QueryInstance) append(searchitem *FindValueOffsetItem, valuedatas []byte, SaveValueSegmentOffset int64) (ValueSegmentOffset uint32, err error) {
	// write data
	segmentOffset := uint32(0)
	if valuedatas != nil {
		segmentOffset, err = ins.writeValueDataToFileWithGC(searchitem, valuedatas)
		if err != nil {
			return 0, err
		}
	} else {
		// do not really write file
		segmentOffset = uint32(SaveValueSegmentOffset)
	}
	// check index type
	ty := searchitem.Type
	if ty == IndexItemTypeValue {
		// value or keep delete mark
		return ins.insertIndexBranch(searchitem, segmentOffset)
	} else if ty == IndexItemTypeNull || ty == IndexItemTypeValueDelete {
		updateitem := searchitem.IncompleteCopy()
		updateitem.Type = IndexItemTypeValue
		updateitem.ValueSegmentOffset = segmentOffset
		return ins.updateSearchItem(updateitem)
	} else {
		return 0, fmt.Errorf("searchitem.Type error, index file breakdown.")
	}
}

/**
 * insert index branch
 */
func (ins *QueryInstance) insertIndexBranch(searchitem *FindValueOffsetItem, valueSegmentOffset uint32) (ValueSegmentOffset uint32, err error) {
	doSaveSearchHash := ins.searchHash
	existSearchHash, _, _ := ins.db.spreadHashToIndexPath(ins.db.convertKeyToHash(searchitem.ValueKey))
	search_i := searchitem.searchCount // already drop file part prefix from hash
	//fmt.Println("search_i", search_i, "key", doSaveSearchHash, existSearchHash)

	idxstat, e1 := ins.targetFilePackage.indexFile.Stat()
	if e1 != nil {
		return 0, e1
	}
	wtatsz := idxstat.Size()
	indexoldsngofst := wtatsz / int64(IndexMenuSize)
	indexcursngofst := indexoldsngofst
	indexFileAppendBytes := bytes.NewBuffer([]byte{})
	for {
		search_i++
		if search_i >= len(doSaveSearchHash) {
			return 0, fmt.Errorf("overflow search hash length.")
		}
		char_1 := doSaveSearchHash[search_i]
		char_2 := existSearchHash[search_i]
		if char_1 == char_2 {
			indexcursngofst++
			brhdts := ins.parseSearchMenu(int(char_2), IndexItemTypeBranch, uint32(indexcursngofst), -1, 0, 0)
			indexFileAppendBytes.Write(brhdts)
			// fmt.Println("- - - Branch brhdts", brhdts)
		} else {
			// can keep delete mark
			searchItemType := searchitem.Type
			if searchItemType != IndexItemTypeValueDelete {
				searchItemType = IndexItemTypeValue
			}
			// write
			brhdts := ins.parseSearchMenu(int(char_1), IndexItemTypeValue, valueSegmentOffset, int(char_2), searchItemType, searchitem.ValueSegmentOffset)
			indexFileAppendBytes.Write(brhdts)
			// fmt.Println("@ value brhdts", brhdts)
			break
		}
	}
	// append index
	appendidxcon := indexFileAppendBytes.Bytes()
	wn, e2 := ins.targetFilePackage.indexFile.WriteAt(appendidxcon, wtatsz)
	if e2 != nil {
		return 0, e2
	}
	if wn != len(appendidxcon) {
		return 0, fmt.Errorf("write to index file error.")
	}

	// update ptr
	brhwriteitem := searchitem.IncompleteCopy()
	brhwriteitem.Type = IndexItemTypeBranch
	brhwriteitem.ValueSegmentOffset = uint32(indexoldsngofst) // start pos
	return ins.updateSearchItem(brhwriteitem)
}

/**
 * insert index branch
 *
func (ins *QueryInstance) insertIndexBranch_old(searchitem *FindValueOffsetItem, valueSegmentOffset uint32) (error) {
	doSaveSearchHash := ins.searchHash
	existSearchHash, _, _ := ins.db.spreadHashToIndexPath( ins.db.convertKeyToHash(searchitem.ValueKey) )
	search_i := searchitem.searchCount // already drop file part prefix from hash
	fmt.Println("search_i", search_i, "key", doSaveSearchHash, existSearchHash)
	canToAddBranch := true
	for {
		search_i ++
		if search_i >= len(doSaveSearchHash) {
			return fmt.Errorf("overflow search hash length.")
		}
		char_1 := doSaveSearchHash[search_i]
		char_2 := existSearchHash[search_i]
		brhwriteitem := searchitem.IncompleteCopy()
		if char_1 != char_2 {
			// ok end
			idxfilesz, e2 := ins.appendSearchMenu(int(char_1), IndexItemTypeValue, valueSegmentOffset, int(char_2), IndexItemTypeValue, searchitem.ValueSegmentOffset)
			if e2 != nil {
				return e2
			}
			brhwriteitem.Type = IndexItemTypeBranch
			brhwriteitem.ValueSegmentOffset = uint32(idxfilesz/int64(IndexMenuSize)) - 1 // start pos
			return ins.updateSearchItem(brhwriteitem)
		}
		// continue branch
		if canToAddBranch {
			idxfilesz, e4 := ins.appendSearchMenu(int(char_2), IndexItemTypeValue, searchitem.ValueSegmentOffset, -1, 0, 0)
			if e4 != nil {
				return e4
			}
			brhwriteitem.Type = IndexItemTypeBranch
			brhwriteitem.ValueSegmentOffset = uint32(idxfilesz/int64(IndexMenuSize)) - 1 // start pos
			e5 := ins.updateSearchItem(brhwriteitem)
			if e5 != nil {
				return e5
			}
			// update old
			//searchitem.Type = IndexItemTypeValue
			//searchitem.ValueSegmentOffset = searchitem.ValueSegmentOffset
			searchitem.IndexMenuSelfSegmentOffset = brhwriteitem.ValueSegmentOffset

		}
		canToAddBranch = true
		//searchitem.IndexItemSelfAlignment = uint32(char_2) * uint32(IndexItemSize)
		continue
	}
}
*/
