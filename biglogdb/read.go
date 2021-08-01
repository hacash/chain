package biglogdb

import "fmt"

// Read log all data
func (db *BigLogDB) Read(key []byte, readlen uint32) ([]byte, error) {
	loghead, logptritem, err := db.ReadHead(key)
	if err != nil {
		return nil, err
	}
	if loghead == nil {
		return nil, nil // not find
	}
	logbody, e1 := db.ReadBodyByPosition(logptritem, readlen)
	if e1 != nil {
		return nil, e1
	}
	return append(loghead, logbody...), nil
}

// Read log Head
func (db *BigLogDB) ReadHead(key []byte) ([]byte, *LogFilePtrSeek, error) {
	// query
	query, err := db.basedb.CreateNewQueryInstance(key)
	if err != nil {
		return nil, nil, err
	}
	defer query.Destroy()
	// read
	vdatas, e1 := query.Find()
	if e1 != nil {
		return nil, nil, e1
	}
	if vdatas == nil {
		return nil, nil, nil // not find
	}
	if len(vdatas) < db.config.LogHeadMaxSize+LogFilePtrSeekSize {
		return nil, nil, fmt.Errorf("db.config.LogHeadMaxSize or stoitem data error.")
	}
	var stoptritem LogFilePtrSeek
	_, e2 := stoptritem.Parse(vdatas, 0)
	if e2 != nil {
		return nil, nil, e2
	}
	// return
	return vdatas[LogFilePtrSeekSize:], &stoptritem, nil
}

// Read data by position
func (db *BigLogDB) ReadBodyByPosition(ptrseek *LogFilePtrSeek, readlen uint32) ([]byte, error) {
	if readlen == 0 || readlen > ptrseek.Valsize {
		readlen = ptrseek.Valsize
	}
	stofile, e0 := db.getStoreFileByNum(ptrseek.Filenum)
	if e0 != nil {
		return nil, e0
	}
	datas := make([]byte, readlen)
	rn, e1 := stofile.ReadAt(datas, int64(ptrseek.Fileseek))
	if e1 != nil {
		return nil, e1
	}
	if rn != int(readlen) {
		return nil, fmt.Errorf("read part store file error.")
	}
	// ret ok
	return datas, nil
}

/////////////////////////////////////////////
