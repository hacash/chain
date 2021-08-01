package blockstorev2

import "github.com/hacash/core/fields"

// trs is exist
func (cs *BlockStore) DeleteTransactionByHash(txhash fields.Hash) error {
	query, e1 := cs.trsdataptrDB.CreateNewQueryInstance(txhash)
	if e1 != nil {
		return e1
	}
	defer query.Destroy()
	return query.Delete()
}
