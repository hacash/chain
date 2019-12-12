package blockstore

import (
	"fmt"
	"github.com/hacash/core/account"
	"github.com/hacash/core/actions"
	"github.com/hacash/core/blocks"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/genesis"
	"github.com/hacash/core/transactions"
	"testing"
)

func Test_t1(t *testing.T) {

	testdir := "/home/shiqiujie/Desktop/Hacash/go/src/github.com/hacash/chain/chainstore/testdata"

	//os.RemoveAll(testdir)

	cscnf := NewEmptyBlockStoreConfig()
	cscnf.Datadir = testdir
	db, e1 := NewBlockStore(cscnf)
	if e1 != nil {
		fmt.Println(e1)
		return
	}

	// create account
	account1 := account.CreateAccountByPassword("123456")
	account2 := account.CreateAccountByPassword("qwerty")
	amount1 := fields.NewAmountNumSmallCoin(12)
	// create action
	action1 := actions.NewAction_1_SimpleTransfer(account2.Address, amount1)
	// create tx
	tx1, e2 := transactions.NewEmptyTransaction_2_Simple(account1.Address)
	tx1.Timestamp = 456789
	if e2 != nil {
		fmt.Println(e2)
		return
	}
	tx1.AppendAction(action1)
	// create block
	block1 := blocks.NewEmptyBlock_v1(genesis.GetGenesisBlock())
	block1.Timestamp = 123456
	coinbasetx1 := transactions.NewTransaction_0_Coinbase()
	coinbasetx1.Address = account1.Address
	coinbasetx1.Reward = *(fields.NewAmountNumSmallCoin(1))
	coinbasetx1.Message = "firstcoinbase"
	block1.AddTransaction(coinbasetx1)
	block1.AddTransaction(tx1)
	block1.SetMrklRoot(blocks.CalculateMrklRoot(block1.GetTransactions())) // update mrkl root
	// serialize
	blkdatas, e3 := block1.Serialize()
	if e3 != nil {
		fmt.Println(e3)
		return
	}

	fmt.Println(len(blkdatas), db)

	fmt.Println(tx1.HashFresh())
	fmt.Println(tx1.Serialize())

	fmt.Println(db.ReadTransactionDataByHash(tx1.HashFresh()))

	//fmt.Println(block1.HashFresh())
	//fmt.Println(0, blkdatas)
	// save
	//err4 := db.SaveBlockUniteTransactions(block1)
	//if err4 != nil {
	//	fmt.Println(err4)
	//	return
	//}

	//fmt.Println(db)

	//readblockdata, e4 := db.ReadBlockDataByHash( block1.HashFresh() )
	//if e4 != nil {
	//	fmt.Println(e4)
	//	return
	//}
	//fmt.Println(1, readblockdata)
	//readblockdata, e4 = db.ReadBlockDataByHeight( block1.GetHeight() )
	//if e4 != nil {
	//	fmt.Println(e4)
	//	return
	//}
	//fmt.Println(2, readblockdata)

}
