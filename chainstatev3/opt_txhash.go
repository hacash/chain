package chainstatev3

import (
	"github.com/hacash/core/fields"
)

// 检查交易是否上链
func (cs *ChainState) CheckTxHash(txhx fields.Hash) (bool, error) {

	value, ok, e := cs.find(KeySuffixType_txhxchk, txhx)
	if e != nil {
		return false, e
	}
	if !ok {
		return false, nil // not find
	}
	if len(value) == 1 && value[0] == 1 {
		return true, nil
	}
	return false, nil
}

// 写入包含交易哈希
func (cs *ChainState) ContainTxHash(txhx fields.Hash) error {
	return cs.save(KeySuffixType_txhxchk, txhx, []byte{1})
}

// 移除交易
func (cs *ChainState) RemoveTxHash(txhx fields.Hash) error {
	return cs.delete(KeySuffixType_txhxchk, txhx)
}
