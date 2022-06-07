package chainstatev3

import (
	"github.com/hacash/core/blocks"
	"github.com/hacash/core/fields"
)

// Read transaction content
func (cs *ChainState) ReadTransactionBytesByHash(txhx fields.Hash) (fields.BlockHeight, []byte, error) {
	hei, e := cs.ReadTxBelongHeightByHash(txhx)
	if e != nil {
		return 0, nil, e
	}
	if hei == 0 {
		return 0, nil, nil // notfind
	}
	_, blkbts, e := cs.BlockStore().ReadBlockBytesByHeight(uint64(hei))
	if e != nil {
		return 0, nil, e
	}
	if blkbts == nil {
		return 0, nil, nil // notfind
	}
	blkObj, _, e := blocks.ParseBlock(blkbts, 0)
	if e != nil {
		return 0, nil, e
	}
	for _, t := range blkObj.GetTrsList() {
		if t.Hash().Equal(txhx) {
			txbtx, e := t.Serialize()
			if e != nil {
				return 0, nil, e
			}
			return hei, txbtx, nil // find success
		}
	}
	return 0, nil, nil // notfind
}

// Check the height of the block to which the exchange belongs
func (cs *ChainState) ReadTxBelongHeightByHash(txhx fields.Hash) (fields.BlockHeight, error) {
	value, ok, e := cs.find(KeySuffixType_txhxchk, txhx)
	if e != nil {
		return 0, e
	}
	if !ok {
		return 0, nil // not find
	}
	if len(value) == int(fields.BlockHeightSize) {
		var hei fields.BlockHeight
		hei.Parse(value, 0)
		return hei, nil // find ok
	}
	return 0, nil
}

// Check whether the transaction is linked
func (cs *ChainState) CheckTxHash(txhx fields.Hash) (bool, error) {
	hei, e := cs.ReadTxBelongHeightByHash(txhx)
	return hei > 0, e
}

// Write include transaction hash
func (cs *ChainState) ContainTxHash(txhx fields.Hash, blkhei fields.BlockHeight) error {
	heibts, _ := blkhei.Serialize()
	return cs.save(KeySuffixType_txhxchk, txhx, heibts)
}

// Remove transaction
func (cs *ChainState) RemoveTxHash(txhx fields.Hash) error {
	return cs.delete(KeySuffixType_txhxchk, txhx)
}
