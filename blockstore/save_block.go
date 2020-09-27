package blockstore

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/hacash/core/blocks"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
)

type transactionStoreItemV1 struct {
	txhash      fields.Hash
	txbody      []byte
	storeoffset uint32
}

// block data store
func (cs *BlockStore) SaveBlockUniteTransactions(fullblock interfaces.Block) error {
	// trs
	allTransactionStoreItem := make([]transactionStoreItemV1, 0, fullblock.GetTransactionCount())
	alltrs := fullblock.GetTransactions()
	if len(alltrs) < 1 {
		return fmt.Errorf("not find coinbase trs.")
	}
	blockdata := bytes.NewBuffer([]byte{})
	// block
	blkheaddata, e1 := fullblock.SerializeHead()
	if e1 != nil {
		return e1
	}
	if len(blkheaddata) != blocks.BlockHeadSize {
		return fmt.Errorf("len(blkheaddata) != blocks.BlockHeadSize, need: %d but got: %d", blocks.BlockHeadSize, len(blkheaddata))
	}
	blkmetadata, e2 := fullblock.SerializeMeta()
	if e2 != nil {
		return e2
	}
	coinbasetrsdata, e3 := alltrs[0].Serialize()
	if e3 != nil {
		return e3
	}
	logprefixofstlen := uint32(len(blkmetadata) + len(coinbasetrsdata))
	// parse trs
	blockdata.Write(blkheaddata)
	blockdata.Write(blkmetadata)
	blockdata.Write(coinbasetrsdata)
	for i := 1; i < len(alltrs); i++ {
		curtrsdata, e := alltrs[i].Serialize()
		if e != nil {
			return e
		}
		allTransactionStoreItem = append(allTransactionStoreItem, transactionStoreItemV1{
			txhash:      alltrs[i].Hash(),
			txbody:      curtrsdata,
			storeoffset: logprefixofstlen,
		})
		logprefixofstlen += uint32(len(curtrsdata))
		blockdata.Write(curtrsdata) // store
	}
	fullblockhash := fullblock.HashFresh()
	// do save block
	stoptr, e5 := cs.blockdataDB.Save(fullblockhash, blockdata.Bytes())
	if e5 != nil {
		return e5
	}
	// do dave blk num
	blknumdata := make([]byte, 8)
	binary.BigEndian.PutUint64(blknumdata, fullblock.GetHeight())
	blkhsquery, e5 := cs.blknumhashDB.CreateNewQueryInstance(blknumdata)
	if e5 != nil {
		return e5
	}
	_, e6 := blkhsquery.Save(fullblockhash)
	if e6 != nil {
		return e6
	}
	blkhsquery.Destroy() // clean
	// do save all trs
	//fmt.Println(stoptr.Filenum, stoptr.Fileseek, stoptr.Valsize)
	for _, item := range allTransactionStoreItem {
		trssto := stoptr.Copy()
		trssto.Fileseek += item.storeoffset
		trssto.Valsize = uint32(len(item.txbody))
		query, e1 := cs.trsdataptrDB.CreateNewQueryInstance(item.txhash)
		if e1 != nil {
			return e1
		}
		trsdataptr, e2 := trssto.Serialize()
		if e2 != nil {
			return e2
		}
		blkhei := fields.VarUint5(fullblock.GetHeight())
		blkheibts, e3 := blkhei.Serialize()
		if e3 != nil {
			return e3
		}
		//fmt.Println(trssto.Filenum, trssto.Fileseek, trssto.Valsize)
		_, e4 := query.Save(append(blkheibts, trsdataptr...))
		if e4 != nil {
			return e4
		}
		query.Destroy() // clean
	}
	return nil
}
