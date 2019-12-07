package hashtreedb

/**
 *
 */
func (ins *QueryInstance) replace(searchitem *FindValueOffsetItem, valuedatas []byte, SaveValueSegmentOffset uint32) (ValueSegmentOffset uint32, err error) {
	if valuedatas != nil {
		return ins.writeSegmentData(searchitem.ValueSegmentOffset, valuedatas)
	} else {
		return SaveValueSegmentOffset, nil // do not really write
	}
}
