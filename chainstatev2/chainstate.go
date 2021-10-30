package chainstatev2

import (
	"fmt"
	"github.com/hacash/chain/leveldb"
	"github.com/hacash/chain/statedomaindb"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/stores"
	"strings"
	"sync"
	"time"
)

type ChainState struct {
	//temporaryDataDir string

	// parent state
	base *ChainState

	// config
	config *ChainStateConfig

	// level db
	ldb *leveldb.DB

	/////////////////////////////

	// status
	laststatusDB *statedomaindb.StateDomainDB

	// state
	balanceDB *statedomaindb.StateDomainDB
	diamondDB *statedomaindb.StateDomainDB
	channelDB *statedomaindb.StateDomainDB // key len = 16
	movebtcDB *statedomaindb.StateDomainDB // 转移BTC记录
	lockblsDB *statedomaindb.StateDomainDB // 线性锁仓地址 key len = 18
	dmdlendDB *statedomaindb.StateDomainDB // 钻石系统借贷 key len = 14
	btclendDB *statedomaindb.StateDomainDB // 钻石系统借贷 key len = 15
	usrlendDB *statedomaindb.StateDomainDB // 钻石系统借贷 key len = 17
	chaswapDB *statedomaindb.StateDomainDB // channel 和 chain 原子互换 key len = 16

	// store
	datastore          interfaces.BlockStore
	datastore_mycreate interfaces.BlockStore

	// data hold
	pendingBlockHeight *uint64
	pendingBlockHash   fields.Hash

	submitStoreDiamond *stores.DiamondSmelt

	prev288BlockTimestamp           uint64
	lastestBlockHeadAndMeta         interfaces.Block
	lastestDiamond                  *stores.DiamondSmelt
	lastestBlockHeadAndMeta_forSave interfaces.Block
	lastestDiamond_forSave          *stores.DiamondSmelt
	totalSupply                     *stores.TotalSupply
	totalSupply_forSave             *stores.TotalSupply

	// status
	isInTxPool bool

	// Read/Write Mutex
	chainStateMutex sync.RWMutex
}

func NewChainState(cnf *ChainStateConfig) (*ChainState, error) {
	return newChainStateEx(cnf, false)
}

func newChainStateEx(cnf *ChainStateConfig, isSubBranchTemporary bool) (*ChainState, error) {
	var e error
	var useldb *leveldb.DB = nil

	if !isSubBranchTemporary {
		useldb, e = leveldb.OpenFile(cnf.Datadir, nil)
		if e != nil {
			return nil, e
		}
	}

	// laststatusDB
	lstcnf := statedomaindb.NewStateDomainDBConfig("laststatus", 0, 0)
	if isSubBranchTemporary {
		lstcnf.MemoryStorage = true // 内存数据库
	} else {
		lstcnf.LevelDB = true // 使用 level db
	}
	laststatusDB := statedomaindb.NewStateDomainDB(lstcnf, useldb)

	// balanceDB
	blscnf := statedomaindb.NewStateDomainDBConfig("balance", 0, fields.AddressSize)
	if isSubBranchTemporary {
		blscnf.MemoryStorage = true // 内存数据库
	} else {
		blscnf.LevelDB = true // 使用 level db
	}
	balanceDB := statedomaindb.NewStateDomainDB(blscnf, useldb)

	/////////////  TEST BEGIN /////////////
	//go func() {
	//	if !isSubBranchTemporary {
	//		time.Sleep(time.Second * 5)
	//		Test_print_all_address_balance(balanceDB)
	//	}
	//}()
	/////////////  TEST END   /////////////

	// diamondDB
	dmdcnf := statedomaindb.NewStateDomainDBConfig("diamond", 0, fields.DiamondNameSize)
	if isSubBranchTemporary {
		dmdcnf.MemoryStorage = true // 内存数据库
	} else {
		dmdcnf.LevelDB = true // 使用 level db
	}
	diamondDB := statedomaindb.NewStateDomainDB(dmdcnf, useldb)

	// channelDB
	chlcnf := statedomaindb.NewStateDomainDBConfig("channel", 0, stores.ChannelIdLength)
	if isSubBranchTemporary {
		chlcnf.MemoryStorage = true // 内存数据库
	} else {
		chlcnf.LevelDB = true // 使用 level db
	}
	channelDB := statedomaindb.NewStateDomainDB(chlcnf, useldb)

	// movebtcDB
	mvbtcnf := statedomaindb.NewStateDomainDBConfig("movebtc", 32, 4)
	if isSubBranchTemporary {
		mvbtcnf.MemoryStorage = true // 内存数据库
	} else {
		mvbtcnf.LevelDB = true // 使用 level db
	}
	movebtcDB := statedomaindb.NewStateDomainDB(mvbtcnf, useldb)

	// lockblsDB
	lkblscnf := statedomaindb.NewStateDomainDBConfig("lockbls", 0, stores.LockblsIdLength)
	//blscnf.KeyReverse = true // 倒排key
	if isSubBranchTemporary {
		lkblscnf.MemoryStorage = true // 内存数据库
	} else {
		lkblscnf.LevelDB = true // 使用 level db
	}
	lockblsDB := statedomaindb.NewStateDomainDB(lkblscnf, useldb)

	// dmdlendDB
	dmdlendcnf := statedomaindb.NewStateDomainDBConfig("dmdlend", 0, stores.DiamondSyslendIdLength)
	if isSubBranchTemporary {
		dmdlendcnf.MemoryStorage = true // 内存数据库
	} else {
		dmdlendcnf.LevelDB = true // 使用 level db
	}
	dmdlendDB := statedomaindb.NewStateDomainDB(dmdlendcnf, useldb)

	// btclendDB
	btclendcnf := statedomaindb.NewStateDomainDBConfig("btclend", 0, stores.BitcoinSyslendIdLength)
	if isSubBranchTemporary {
		btclendcnf.MemoryStorage = true // 内存数据库
	} else {
		btclendcnf.LevelDB = true // 使用 level db
	}
	btclendDB := statedomaindb.NewStateDomainDB(btclendcnf, useldb)

	// btclendDB
	usrlendcnf := statedomaindb.NewStateDomainDBConfig("usrlend", 0, stores.UserLendingIdLength)
	if isSubBranchTemporary {
		usrlendcnf.MemoryStorage = true // 内存数据库
	} else {
		usrlendcnf.LevelDB = true // 使用 level db
	}
	usrlendDB := statedomaindb.NewStateDomainDB(usrlendcnf, useldb)

	// chaswapDB
	chaswapcnf := statedomaindb.NewStateDomainDBConfig("chaswap", 0, fields.HashHalfCheckerSize)
	if isSubBranchTemporary {
		chaswapcnf.MemoryStorage = true // 内存数据库
	} else {
		chaswapcnf.LevelDB = true // 使用 level db
	}
	chaswapDB := statedomaindb.NewStateDomainDB(chaswapcnf, useldb)

	// return ok
	cs := &ChainState{
		config: cnf,
		// temporaryDataDir:      temporaryDataDir,
		base:                  nil,
		ldb:                   useldb,
		laststatusDB:          laststatusDB,
		balanceDB:             balanceDB,
		diamondDB:             diamondDB,
		channelDB:             channelDB,
		movebtcDB:             movebtcDB,
		lockblsDB:             lockblsDB,
		dmdlendDB:             dmdlendDB,
		btclendDB:             btclendDB,
		usrlendDB:             usrlendDB,
		chaswapDB:             chaswapDB,
		prev288BlockTimestamp: 0,
		pendingBlockHeight:    nil,
		pendingBlockHash:      nil,
		isInTxPool:            false,
		chainStateMutex:       sync.RWMutex{},
	}

	/////////////  TEST BEGIN /////////////
	//go func() {
	//	if !isSubBranchTemporary {
	//		time.Sleep(time.Second * 5)
	//		Test_print_283756384756389(cs)
	//	}
	//}()
	/////////////  TEST END   /////////////

	return cs, nil
}

// interface api
func (cs *ChainState) Destory() {
	cs.DestoryTemporary()
}

func (cs *ChainState) Close() {
	if cs.ldb != nil {
		cs.ldb.Close()
		cs.ldb = nil
	}
	if cs.laststatusDB != nil {
		cs.laststatusDB.Close()
	}
	if cs.balanceDB != nil {
		cs.balanceDB.Close()
	}
	if cs.channelDB != nil {
		cs.channelDB.Close()
	}
	if cs.diamondDB != nil {
		cs.diamondDB.Close()
	}
	if cs.movebtcDB != nil {
		cs.movebtcDB.Close()
	}
	if cs.lockblsDB != nil {
		cs.lockblsDB.Close()
	}
	if cs.dmdlendDB != nil {
		cs.dmdlendDB.Close()
	}
	if cs.btclendDB != nil {
		cs.btclendDB.Close()
	}
	if cs.usrlendDB != nil {
		cs.usrlendDB.Close()
	}
	if cs.chaswapDB != nil {
		cs.chaswapDB.Close()
	}
	// close mycreate store
	if cs.datastore_mycreate != nil {
		cs.datastore_mycreate.Close()
	}
}

// destory temporary
func (cs *ChainState) DestoryTemporary() {
}

// chain data store
func (cs *ChainState) BlockStore() interfaces.BlockStore {
	if cs.datastore != nil {
		return cs.datastore
	}
	if cs.base != nil {
		cs.datastore = cs.base.BlockStore() // copy
		return cs.datastore
	}
	return cs.datastore
}

func (cs *ChainState) SetBlockStore(store interfaces.BlockStore) error {
	if cs.base != nil {
		return fmt.Errorf("Can only be set chainstore in the final state.")
	}
	cs.datastore = store
	cs.datastore_mycreate = store // 我创建的
	return nil
}

// merge write
func (cs *ChainState) MergeCoverWriteChainState(src *ChainState) error {
	// copy status
	if src.pendingBlockHeight != nil {
		cs.SetPendingBlockHeight(*src.pendingBlockHeight)
	}
	if src.pendingBlockHash != nil {
		cs.SetPendingBlockHash(src.pendingBlockHash)
	}
	if src.submitStoreDiamond != nil {
		e := cs.SetPendingSubmitStoreDiamond(src.submitStoreDiamond)
		if e != nil {
			return e
		}
	}
	if src.lastestBlockHeadAndMeta != nil {
		e := cs.SetLastestBlockHeadAndMeta(src.lastestBlockHeadAndMeta)
		if e != nil {
			return e
		}
	}
	if src.lastestDiamond != nil {
		e := cs.SetLastestDiamond(src.lastestDiamond)
		if e != nil {
			return e
		}
	}
	if src.totalSupply != nil {
		e := cs.UpdateSetTotalSupply(src.totalSupply)
		if e != nil {
			return e
		}
	}

	//  COPY COVER WRITE STATE

	e1 := cs.balanceDB.TraversalCopy(src.balanceDB)
	if e1 != nil {
		return e1
	}
	e2 := cs.diamondDB.TraversalCopy(src.diamondDB)
	if e2 != nil {
		return e2
	}
	e3 := cs.channelDB.TraversalCopy(src.channelDB)
	if e3 != nil {
		return e3
	}
	e5 := cs.movebtcDB.TraversalCopy(src.movebtcDB)
	if e5 != nil {
		return e5
	}
	e6 := cs.lockblsDB.TraversalCopy(src.lockblsDB)
	if e6 != nil {
		return e6
	}
	e7 := cs.dmdlendDB.TraversalCopy(src.dmdlendDB)
	if e7 != nil {
		return e7
	}
	e8 := cs.btclendDB.TraversalCopy(src.btclendDB)
	if e8 != nil {
		return e8
	}
	e9 := cs.usrlendDB.TraversalCopy(src.usrlendDB)
	if e9 != nil {
		return e9
	}
	e10 := cs.chaswapDB.TraversalCopy(src.chaswapDB)
	if e10 != nil {
		return e10
	}

	// copy ok

	return nil
}

// interface api
func (cs *ChainState) Fork() (interfaces.ChainState, error) {
	return cs.NewSubBranchTemporaryChainState()
}

func (cs *ChainState) IsDatabaseVersionRebuildMode() bool {
	if cs.base != nil {
		// 递归向上
		return cs.base.IsDatabaseVersionRebuildMode()
	}
	// 返回最终配置
	return cs.config.DatabaseVersionRebuildMode
}

// fork sub
func (cs *ChainState) NewSubBranchTemporaryChainState() (*ChainState, error) {

	tempcnf := NewEmptyChainStateConfig()
	// 拷贝一些配置
	tempcnf.BTCMoveCheckEnable = cs.config.BTCMoveCheckEnable
	tempcnf.BTCMoveCheckLogsURL = cs.config.BTCMoveCheckLogsURL
	// 拷贝完毕
	newTempState, err1 := newChainStateEx(tempcnf, true)
	if err1 != nil {
		return nil, err1
	}
	// set base
	newTempState.base = cs
	// ok
	return newTempState, nil
}

// submit to write disk
func (cs *ChainState) SubmitDataStoreWriteToInvariableDisk(block interfaces.Block) error {
	if cs.base != nil {
		panic("Can only be saved in the final state.")
	}
	if cs.pendingBlockHash == nil {
		//return fmt.Errorf("Block pending hash not set.")
		panic("pending block hash not be set.")
	}
	//
	store := cs.BlockStore()
	if store == nil {
		return fmt.Errorf("Not find BlockStore object.")
	}
	// save status
	e0 := cs.SetLastestBlockHeadAndMeta(block)
	if e0 != nil {
		return e0
	}
	e1 := cs.IncompleteSaveLastestBlockHeadAndMeta()
	if e1 != nil {
		return e1
	}
	e2 := cs.IncompleteSaveLastestDiamond()
	if e2 != nil {
		return e2
	}
	e9 := cs.IncompleteSaveTotalSupply()
	if e9 != nil {
		return e9
	}
	// save diamond
	if cs.submitStoreDiamond != nil {
		cs.submitStoreDiamond.ContainBlockHash = cs.pendingBlockHash // copy
		e := store.SaveDiamond(cs.submitStoreDiamond)
		if e != nil {
			return e
		}
	}
	// save block data
	e3 := store.SaveBlockUniteTransactions(block)
	if e3 != nil {
		return e3
	}
	// reset clear
	cs.chainStateMutex.Lock()
	cs.pendingBlockHash = nil
	cs.pendingBlockHeight = nil
	cs.submitStoreDiamond = nil
	cs.chainStateMutex.Unlock()
	// ok
	return nil
}

////////////////////////////////////////////////////

func Test_print_283756384756389(db *ChainState) {

	ttchun := 0

	for i := 1; ; i++ {
		diaobj, e := db.datastore.ReadDiamondByNumber(uint32(i))
		if e != nil {
			break
		}
		vg := diaobj.VisualGene
		if vg[7] == vg[8] {
			hexstr := vg.ToHex()
			if hexstr[17] == hexstr[16] &&
				hexstr[16] == hexstr[15] &&
				hexstr[15] == hexstr[14] {
				ttchun += 1
				fmt.Println(diaobj.Number, string(diaobj.Diamond))
			}
		}
	}

	fmt.Println("total: ", ttchun)

	/*


	YHSBKT,AEZYXZ,UHIWIW,SNSVHB,HVBSMA,XUISAZ,VAEEMZ,XIHKSY,NAHMIT,XWUTXM,ATVBWI,WSKWHT,AHANHU,XAESVH,AXTZZS


	*/

}

func Test_print_all_address_balance(db *statedomaindb.StateDomainDB) {

	time.Sleep(time.Microsecond)

	total_address_count := int64(0)
	total_hac_address_count := int64(0)
	total_btc_address_count := int64(0)
	total_hacd_address_count := int64(0)

	total_hac := float64(0)
	total_btc := int64(0)
	total_hacd := int(0)

	iter := db.LevelDB.NewIterator(nil, nil)
	for iter.Next() {
		//fmt.Printf("key:%s, value:%s\n", iter.Key(), iter.Value())
		key := iter.Key()
		if !strings.HasSuffix(string(key), "balance") {
			continue
		}
		addrbt := iter.Key()[0:21]
		addr := fields.Address(addrbt)
		var bls = stores.Balance{}
		bls.Parse(iter.Value(), 0)
		hacfltn := bls.Hacash.ToMei()
		if hacfltn == 0 && bls.Satoshi == 0 && bls.Diamond == 0 {
			continue
		}
		fmt.Printf("% 4d %-34s % 12.4f % 6d %d\n", total_address_count+1, addr.ToReadable(), hacfltn, bls.Diamond, bls.Satoshi)
		total_hac += float64(hacfltn)
		total_btc += int64(bls.Satoshi)
		total_hacd += int(bls.Diamond)
		if float64(hacfltn) > 0 {
			total_hac_address_count++
		}
		if int64(bls.Satoshi) > 0 {
			total_btc_address_count++
		}
		if int(bls.Diamond) > 0 {
			total_hacd_address_count++
		}
		total_address_count++

	}
	iter.Release()

	fmt.Printf("------------------\n[ADDRESS] %d address, hac: %d, btc: %d, hacd: %d \n[AMOUNT] HAC: %f, SAT: %d, HACD: %d\n", total_address_count, total_hac_address_count, total_btc_address_count, total_hacd_address_count, total_hac, total_btc, total_hacd)
}
