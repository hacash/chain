package blockstore

import (
	"fmt"
	"github.com/hacash/chain/biglogdb"
	"github.com/hacash/core/fields"
)

// trs is exist
func (cs *BlockStore) TransactionIsExist(txhash fields.Hash) (bool, error) {
	query, e1 := cs.trsdataptrDB.CreateNewQueryInstance(txhash)
	if e1 != nil {
		return false, e1
	}
	defer query.Destroy()
	return query.Exist() // search
}

// block data store
func (cs *BlockStore) ReadTransactionBytesByHash(txhash fields.Hash) (uint64, []byte, error) {
	query, e1 := cs.trsdataptrDB.CreateNewQueryInstance(txhash)
	if e1 != nil {
		return 0, nil, e1
	}
	defer query.Destroy()
	findbytes, e2 := query.Find()
	if e2 != nil {
		return 0, nil, e2
	}
	if findbytes == nil {
		return 0, nil, nil // not find
	}
	if len(findbytes) < 5 {
		return 0, nil, fmt.Errorf("data size error.")
	}
	trsposptrdts := findbytes[5:]
	if len(trsposptrdts) != biglogdb.LogFilePtrSeekSize {
		return 0, nil, fmt.Errorf("log file ptr seek data length is error.")
	}
	// read from block store
	var ptritem biglogdb.LogFilePtrSeek
	_, e3 := ptritem.Parse(trsposptrdts, 0)
	if e3 != nil {
		return 0, nil, e3
	}
	// read trs
	trsdata, e4 := cs.blockdataDB.ReadBodyByPosition(&ptritem, 0)
	if e4 != nil {
		return 0, nil, e4
	}
	// return ok
	height := fields.BlockHeight(0)
	height.Parse(findbytes[0:5], 0)
	return uint64(height), trsdata, nil
}
