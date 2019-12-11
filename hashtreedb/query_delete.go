package hashtreedb

/**
 * clear search index cache
 */
func (ins *QueryInstance) Delete() error {

	if ins.db.config.KeepDeleteMark {

	}

	if ins.db.config.SaveMarkBeforeValue {

	}

	if ins.db.config.ForbidGC {

	}

	return nil
}
