package hashtreedb

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

/**
 * clear search index cache
 */
func (ins *QueryInstance) Delete() error {

	// 内存数据库
	if ins.db.config.MemoryStorage {
		ins.db.MemoryStorageDB.Delete(ins.key)
		return nil
	}
	// 内存数据库
	if ins.db.config.LevelDB {
		ins.db.LevelDB.Delete(ins.key, nil)
		return nil
	}

	panic("NewHashTreeDB  must use LevelDB!")

	// 文件数据库
	ins.ClearSearchIndexCache()
	ofstItem, err := ins.SearchIndex()
	if err != nil {
		return err // error
	}
	if ofstItem == nil {
		return nil // not find ok
	}
	if ofstItem.Type == IndexItemTypeValueDelete {
		return nil // already deleted
	}
	if ofstItem.Type != IndexItemTypeValue {
		return nil // not value
	}
	e2 := ins.readSegmentDataFillItem(ofstItem, false)
	if e2 != nil {
		return e2 // error
	}
	if bytes.Compare(ins.key, ofstItem.ValueKey) != 0 {
		// read target ok other one
		return nil
	}
	if ins.db.config.KeepDeleteMark {
		ofstItem.Type = IndexItemTypeValueDelete // mark delete
	} else {
		ofstItem.Type = 0 // nothing
	}
	var valueSegmentOffset = ofstItem.ValueSegmentOffset
	if !ins.db.config.ForbidGC && !ins.db.config.KeepDeleteMark {
		e := ins.collecteGarbageSpace(ofstItem) // Collecte Garbage Space
		if e != nil {
			return e
		}
		ofstItem.ValueSegmentOffset = 0
	}
	_, e := ins.updateSearchItem(ofstItem) // update index
	if e != nil {
		return e
	}
	// update value
	if ins.db.config.SaveMarkBeforeValue {
		// write delete mark
		ins.writeSegmentDataEx(valueSegmentOffset, []byte{IndexItemTypeValueDelete})
	}
	return nil
}

// get space
func (ins *QueryInstance) releaseGarbageSpace() (uint32, error) {
	if ins.db.config.ForbidGC {
		return 0, nil // ban
	}
	gcf := ins.targetFilePackage.gcFile
	stat, e := gcf.Stat()
	if e != nil {
		return 0, e
	}
	if stat.Size() == 0 {
		return 0, nil // empty
	}
	if stat.Size()%4 != 0 {
		return 0, fmt.Errorf("gc file is break down.")
	}
	gcptr := make([]byte, 4)
	_, e1 := gcf.ReadAt(gcptr, stat.Size()-4)
	if e1 != nil {
		return 0, e1
	}
	// change
	e2 := gcf.Truncate(stat.Size() - 4)
	if e2 != nil {
		return 0, e2
	}
	// ok
	return binary.BigEndian.Uint32(gcptr), nil
}

// addspace
func (ins *QueryInstance) collecteGarbageSpace(item *FindValueOffsetItem) error {
	if ins.db.config.ForbidGC {
		return nil // ban
	}
	gcf := ins.targetFilePackage.gcFile
	stat, e := gcf.Stat()
	if e != nil {
		return e
	}
	valuesegptr := make([]byte, 4)
	binary.BigEndian.PutUint32(valuesegptr, item.ValueSegmentOffset)
	_, e1 := gcf.WriteAt(valuesegptr, stat.Size())
	if e1 != nil {
		return e1
	}
	return nil
}
