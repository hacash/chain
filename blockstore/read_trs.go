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
func (cs *BlockStore) ReadTransactionDataByHash(txhash fields.Hash) ([]byte, error) {
	query, e1 := cs.trsdataptrDB.CreateNewQueryInstance(txhash)
	if e1 != nil {
		return nil, e1
	}
	defer query.Destroy()
	trsposptrdts, e2 := query.Find()
	if e2 != nil {
		return nil, e2
	}
	if trsposptrdts == nil {
		return nil, nil // not find
	}
	if len(trsposptrdts) != biglogdb.LogFilePtrSeekSize {
		return nil, fmt.Errorf("log file ptr seek data length is error.")
	}
	// read from block store
	var ptritem biglogdb.LogFilePtrSeek
	_, e3 := ptritem.Parse(trsposptrdts, 0)
	if e3 != nil {
		return nil, e3
	}
	// read trs
	trsdata, e4 := cs.blockdataDB.ReadBodyByPosition(&ptritem)
	if e4 != nil {
		return nil, e4
	}
	// return ok
	return trsdata, nil
}
