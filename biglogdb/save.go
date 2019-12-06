package biglogdb

import "fmt"

// save data
func (db *BigLogDB) Save(key []byte, valuedata []byte) (*LogFilePtrSeek, error) {
	// check file num
	db.wlock.Lock()
	defer db.wlock.Unlock()
	// if log head
	logheaddatas := []byte{}
	logbodydatas := valuedata
	lghdmaxsz := db.config.LogHeadMaxSize
	if lghdmaxsz > 0 {
		if len(valuedata) < lghdmaxsz {
			return nil, fmt.Errorf("value data length cannot less than db.config.LogHeadMaxSize.")
		}
		logheaddatas = valuedata[0:lghdmaxsz]
		logbodydatas = valuedata[lghdmaxsz:]
	}
	// save to store file
	currentfilenum, e1 := db.GetFileNum()
	if e1 != nil {
		return nil, e1
	}
OPENSTOREFILE:
	stofile, e2 := db.getStoreFileByNum(currentfilenum)
	if e2 != nil {
		return nil, e2
	}
	stofstat, e3 := stofile.Stat()
	if e3 != nil {
		return nil, e3
	}
	stofsize := stofstat.Size()
	if stofsize+int64(len(logbodydatas)) > db.config.BlockPartFileMaxSize {
		if stofsize == 0 && int64(len(logbodydatas)) > db.config.BlockPartFileMaxSize {
			// force use current file
		} else {
			// open next file
			currentfilenum = currentfilenum + 1
			e := db.SetFileNum(currentfilenum)
			if e != nil {
				return nil, e
			}
			goto OPENSTOREFILE // use next file
		}
	}
	// write file to
	wn, e4 := stofile.WriteAt(logbodydatas, stofsize)
	if e4 != nil {
		return nil, e4
	}
	if wn != len(logbodydatas) {
		return nil, fmt.Errorf("store file write error.")
	}
	// ptr seek
	stoptr := &LogFilePtrSeek{
		Filenum:  currentfilenum,
		Fileseek: uint32(stofsize),
		Valsize:  uint32(len(logbodydatas)),
	}
	// save index
	query, err := db.bashhashtreedb.CreateNewQueryInstance(key)
	if err != nil {
		return nil, err
	}
	defer query.Destroy()
	stoptrdts, _ := stoptr.Serialize()
	stoptrdts = append(stoptrdts, logheaddatas...)
	e5 := query.Save(stoptrdts)
	if e5 != nil {
		return nil, e5
	}
	// return ok stoptr
	return stoptr, nil
}
