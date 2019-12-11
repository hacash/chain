package hashtreedb

/**
 *
 */
func (ins *QueryInstance) replace(searchitem *FindValueOffsetItem, valuedatas []byte, SaveValueSegmentOffset int64) (ValueSegmentOffset uint32, err error) {
	if valuedatas != nil {
		return ins.writeSegmentData(searchitem.ValueSegmentOffset, valuedatas)
	} else {
		return uint32(SaveValueSegmentOffset), nil // do not really write
	}
}
