package hashtreedb

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

/**
 * Save Value Segment Offset
 */
func (ins *QueryInstance) UnsafeSaveWithValueSegmentOffset(ValueSegmentOffset uint32) error {
	_, err := ins.saveEx(nil, int64(ValueSegmentOffset))
	return err
}

/**
 * Save Value
 */
func (ins *QueryInstance) Save(valuedatas []byte) (ValueSegmentOffset uint32, err error) {
	return ins.saveEx(valuedatas, -1)
}

/**
 * Save Value
 */
func (ins *QueryInstance) saveEx(valuedatas []byte, SaveValueSegmentOffset int64) (ValueSegmentOffset uint32, err error) {
	if valuedatas == nil && SaveValueSegmentOffset < 0 {
		return 0, fmt.Errorf("valuedatas or SaveValueSegmentOffset must give one.")
	}
	if valuedatas != nil {
		dtleg := len(valuedatas)
		mxvsz := int(ins.db.config.MaxValueSize)
		if dtleg > mxvsz {
			return 0, fmt.Errorf("value size too much.")
		} else if dtleg < mxvsz {
			// add to Max size
			maxvalue := make([]byte, mxvsz)
			copy(maxvalue, valuedatas)
			valuedatas = maxvalue
		}
	}
	ofstitem, err := ins.SearchIndex()
	if err != nil {
		return 0, err
	}
	if ofstitem == nil {
		return 0, fmt.Errorf("*FindValueOffsetItem not find, index file is breakdown.")
	}
	if ofstitem.Type == IndexItemTypeValue {
		e2 := ins.readSegmentDataFillItem(ofstitem, false)
		if e2 != nil {
			return 0, e2 // error
		}
		if bytes.Compare(ins.key, ofstitem.ValueKey) == 0 {
			// replace target data
			return ins.replace(ofstitem, valuedatas, SaveValueSegmentOffset)
		}
	}
	// add new value
	return ins.append(ofstitem, valuedatas, SaveValueSegmentOffset)
}

/**
 * update search item
 */

func (ins *QueryInstance) updateSearchItem(sitem *FindValueOffsetItem) (ValueSegmentOffset uint32, err error) {
	atpos := sitem.IndexMenuSelfSegmentOffset*uint32(IndexMenuSize) + sitem.IndexItemSelfAlignment
	itbytes := sitem.Serialize()
	wn, err := ins.targetFilePackage.indexFile.WriteAt(itbytes, int64(atpos))
	if err != nil {
		return 0, err
	}
	if wn != len(itbytes) {
		return 0, fmt.Errorf("write to index file error.")
	}
	return sitem.ValueSegmentOffset, nil
}

/**
 * update search item
 */

func (ins *QueryInstance) parseSearchMenu(char1 int, mark1 byte, segofst1 uint32, char2 int, mark2 byte, segofst2 uint32) []byte {
	chars := []int{char1, char2}
	marks := []byte{mark1, mark2}
	segofsts := []uint32{segofst1, segofst2}
	menubytes := bytes.Repeat([]byte{0}, IndexMenuSize)
	for i := 0; i < 2; i++ {
		char := chars[i]
		mark := marks[i]
		segofst := segofsts[i]
		if char < 0 || char > 255 || mark < 0 || mark > 255 {
			break
		}
		// insert
		insert_pos := char * IndexItemSize
		menubytes[insert_pos] = mark
		ofstdts := []byte{0, 0, 0, 0}
		binary.BigEndian.PutUint32(ofstdts, segofst)
		copy(menubytes[insert_pos+1:insert_pos+1+4], ofstdts)
	}
	return menubytes
}

func (ins *QueryInstance) appendSearchMenu(char1 int, mark1 byte, segofst1 uint32, char2 int, mark2 byte, segofst2 uint32) (int64, error) {
	menubytes := ins.parseSearchMenu(char1, mark1, segofst1, char2, mark2, segofst2)
	// write file
	idxstat, e1 := ins.targetFilePackage.indexFile.Stat()
	if e1 != nil {
		return 0, e1
	}
	wtatsz := idxstat.Size()
	wn, e2 := ins.targetFilePackage.indexFile.WriteAt(menubytes, wtatsz)
	if e2 != nil {
		return 0, e2
	}
	if wn != IndexMenuSize {
		return 0, fmt.Errorf("write to index file error.")
	}
	return wtatsz + int64(IndexMenuSize), nil
}

/**
 * Save data
 */
func (ins *QueryInstance) writeSegmentDataEx(segmentOffset uint32, segdatas []byte) (ValueSegmentOffset uint32, err error) {
	// write to the file
	wtpos := int64(segmentOffset * ins.db.config.segmentValueSize)
	wn, e4 := ins.targetFilePackage.dataFile.WriteAt(segdatas, wtpos)
	if e4 != nil {
		return 0, e4
	}
	if wn != len(segdatas) {
		return 0, fmt.Errorf("segment file WriteAt length error.")
	}
	return segmentOffset, nil
}

/**
 * Save data
 */
func (ins *QueryInstance) writeSegmentData(segmentOffset uint32, valuedatas []byte) (ValueSegmentOffset uint32, err error) {
	// write to the file
	var datas = bytes.NewBuffer([]byte{})
	if ins.db.config.SaveMarkBeforeValue {
		datas.WriteByte(byte(2)) // store mark
	}
	datas.Write(ins.key)
	datas.Write(valuedatas)
	// write
	segdatas := datas.Bytes()
	return ins.writeSegmentDataEx(segmentOffset, segdatas)
}

/**
 * Save data
 */
func (ins *QueryInstance) writeValueDataToFileWithGC(searchitem *FindValueOffsetItem, valuedatas []byte) (valueSegmentOffset uint32, err error) {
	var segmentwtat int64 = -1
	if searchitem.Type == IndexItemTypeValueDelete {
		segmentwtat = int64(searchitem.ValueSegmentOffset * ins.db.config.segmentValueSize)
	}
	// check gc
	if segmentwtat == -1 && ins.db.config.ForbidGC == false {
		ptr, e := ins.releaseGarbageSpace()
		if e != nil {
			return 0, e
		}
		if ptr > 0 {
			segmentwtat = int64(ptr * ins.db.config.segmentValueSize)
		}
	}
	// append file tail
	if segmentwtat == -1 {
		dtstat, e1 := ins.targetFilePackage.dataFile.Stat()
		if e1 != nil {
			return 0, e1
		}
		segmentwtat = dtstat.Size()
		if uint32(segmentwtat)%ins.db.config.segmentValueSize != 0 {
			err = fmt.Errorf("segmentwtat(%d) %% ins.db.config.segmentValueSize(%d) != 0, data file is breakdown.", uint32(segmentwtat), ins.db.config.segmentValueSize)
			return
		}
	}
	// write segment data
	segmentOffset := uint32(segmentwtat) / ins.db.config.segmentValueSize
	return ins.writeSegmentData(segmentOffset, valuedatas)
}
