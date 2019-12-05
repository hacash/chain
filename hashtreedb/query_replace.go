package hashtreedb






/**
 *
 */
func (ins *QueryInstance) replace(searchitem *FindValueOffsetItem, valuedatas []byte) (error) {
	return ins.writeSegmentData(searchitem.ValueSegmentOffset, valuedatas)
}


