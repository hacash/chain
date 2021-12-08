package chainstatev3

import (
	"fmt"
	"github.com/hacash/chain/blockstorev3"
	"github.com/hacash/chain/leveldb"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/stores"
	"math/rand"
	"strings"
	"sync"
	"time"
)

type ChainState struct {
	sid uint64 // 唯一标识

	// config
	config *ChainStateConfig

	blockstore *blockstorev3.BlockStore

	// parent & childs state
	base   *ChainState
	childs map[uint64]*ChainState

	// level db
	ldb   *leveldb.DB
	memdb *sync.Map

	isInTxPool bool

	pending         interfaces.PendingStatus
	lastStatusCache interfaces.LatestStatus

	// lock
	statusMux *sync.RWMutex
}

func NewChainStateImmutable(cnf *ChainStateConfig) (*ChainState, error) {
	ins, e := newChainStateEx(cnf, false)
	if e != nil {
		return nil, e
	}

	return ins, nil
}

func newChainStateEx(cnf *ChainStateConfig, isSubBranchTemporary bool) (*ChainState, error) {

	sid := rand.Uint64()
	state := &ChainState{
		sid:             sid,
		config:          cnf,
		blockstore:      nil,
		base:            nil,
		childs:          make(map[uint64]*ChainState, 0),
		pending:         nil,
		lastStatusCache: nil,
		ldb:             nil,
		memdb:           nil,
		isInTxPool:      false,
		statusMux:       &sync.RWMutex{},
	}

	// 建立数据库 db
	if isSubBranchTemporary {
		state.memdb = &sync.Map{}
	} else {
		useldb, e := leveldb.OpenFile(cnf.Datadir, nil)
		if e != nil {
			return nil, e
		}
		state.ldb = useldb
		/////////////  TEST BEGIN /////////////
		//go func() {
		//	if !isSubBranchTemporary {
		//		time.Sleep(time.Second * 5)
		//		Test_print_all_address_balance(useldb)
		//	}
		//}()
		/////////////  TEST END   /////////////

	}

	// ok
	return state, nil
}

func (s ChainState) BlockStore() interfaces.BlockStore {
	s.statusMux.RLock()
	defer s.statusMux.RUnlock()

	if s.blockstore != nil {
		return s.blockstore
	}

	if s.base != nil {
		return s.base.BlockStore()
	}

	return nil
}

func (s ChainState) BlockStoreRead() interfaces.BlockStoreRead {
	return s.blockstore
}

func (s *ChainState) SetBlockStoreObj(store *blockstorev3.BlockStore) {
	s.statusMux.Lock()
	defer s.statusMux.Unlock()

	s.blockstore = store
}

func (s ChainState) IsImmutable() bool {
	s.statusMux.RLock()
	defer s.statusMux.RUnlock()

	if s.ldb != nil {
		return true
	}
	return false
}

func (s *ChainState) ForkSubChildObj() (*ChainState, error) {
	s.statusMux.Lock()
	defer s.statusMux.Unlock()

	// create
	sub, e := newChainStateEx(s.config, true)
	if e != nil {
		return nil, e
	}
	// copy field
	sub.config = s.config
	sub.blockstore = s.blockstore
	sub.isInTxPool = s.isInTxPool
	sub.pending = s.pending
	// set base
	sub.base = s
	// add childs
	s.childs[sub.sid] = sub
	// ok
	return sub, nil
}

func (s *ChainState) ForkSubChild() (interfaces.ChainState, error) {
	obj, e := s.ForkSubChildObj()
	if e != nil {
		return nil, e
	}
	return obj, nil
}

func (s *ChainState) ForkNextBlockObj(hei uint64, hx fields.Hash, blockhead interfaces.Block) (*ChainState, error) {
	// create
	sub, e := s.ForkSubChildObj()
	if e != nil {
		return nil, e
	}

	s.statusMux.Lock()
	defer s.statusMux.Unlock()

	// set block
	sub.pending = NewPendingStatus(hei, hx, blockhead)

	// ok
	return sub, nil
}

func (s *ChainState) ForkNextBlock(hei uint64, hx fields.Hash, blockhead interfaces.Block) (interfaces.ChainState, error) {
	obj, e := s.ForkNextBlockObj(hei, hx, blockhead)
	if e != nil {
		return nil, e
	}
	return obj, nil
}

//// 获取指向的区块
//func (s ChainState) GetReferBlock() (uint64, fields.Hash) {
//	s.statusMux.RLock()
//	defer s.statusMux.RUnlock()
//
//	return s.referBlockHeight, s.referBlockHash
//}

// 获得父级状态
func (s ChainState) GetParentObj() *ChainState {
	s.statusMux.RLock()
	defer s.statusMux.RUnlock()

	return s.base
}
func (s ChainState) GetParent() interfaces.ChainState {
	base := s.GetParentObj()
	return base
}

// 获得所有子状态
func (s *ChainState) GetChildObjs() map[uint64]*ChainState {
	s.statusMux.RLock()
	defer s.statusMux.RUnlock()

	return s.childs
}
func (s *ChainState) GetChilds() map[uint64]interfaces.ChainState {
	s.statusMux.RLock()
	defer s.statusMux.RUnlock()

	var childs = make(map[uint64]interfaces.ChainState)
	for i, v := range s.childs {
		childs[i] = v
	}

	return childs
}

// 获得所有子状态
func (s *ChainState) RemoveChild(child *ChainState) {
	s.statusMux.Lock()
	defer s.statusMux.Unlock()

	delete(s.childs, child.sid)
}

// 销毁，包括删除所有子状态、缓存、状态数据等
func (s *ChainState) Destory() {
	s.statusMux.Lock()
	defer s.statusMux.Unlock()

	if s.base != nil {
		s.base.RemoveChild(s)
	}
	// clean
	s.config = nil
	s.blockstore = nil
	s.pending = nil
	s.lastStatusCache = nil
	s.base = nil
	s.childs = nil
	s.memdb = nil
	s.ldb = nil
}

// 关闭文件句柄等
func (s *ChainState) Close() {
	s.statusMux.Lock()
	defer s.statusMux.Unlock()

	if s.ldb != nil {
		// close
		s.ldb.Close()
		s.ldb = nil
	}
}

func (s *ChainState) IsDatabaseVersionRebuildMode() bool {
	s.statusMux.RLock()
	defer s.statusMux.RUnlock()

	// 返回最终配置
	return s.config.DatabaseVersionRebuildMode
}

func (s *ChainState) SetDatabaseVersionRebuildMode(set bool) {
	s.statusMux.RLock()
	defer s.statusMux.RUnlock()

	s.config.DatabaseVersionRebuildMode = set
}

func (s *ChainState) SetInTxPool(stat bool) {
	s.statusMux.Lock()
	defer s.statusMux.Unlock()

	s.isInTxPool = stat
}

func (s *ChainState) IsInTxPool() bool {
	s.statusMux.RLock()
	defer s.statusMux.RUnlock()

	return s.isInTxPool
}

/////////////////////////////////////////////

func Test_print_all_address_balance(ldb *leveldb.DB) {

	time.Sleep(time.Microsecond)

	total_address_count := int64(0)
	total_hac_address_count := int64(0)
	total_btc_address_count := int64(0)
	total_hacd_address_count := int64(0)

	total_hac := float64(0)
	total_btc := int64(0)
	total_hacd := int(0)

	iter := ldb.NewIterator(nil, nil)
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
