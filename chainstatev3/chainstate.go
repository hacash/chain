package chainstatev3

import (
	"github.com/hacash/chain/leveldb"
	"github.com/hacash/core/fields"
	"math/rand"
	"sync"
)

type ChainState struct {
	sid uint64 // 唯一标识

	// config
	config *ChainStateConfig

	// parent & childs state
	base   *ChainState
	childs []*ChainState

	referBlockHeight uint64
	referBlockHash   fields.Hash

	// level db
	ldb   *leveldb.DB
	memdb *sync.Map

	// lock
	statusMux *sync.RWMutex
}

func NewChainStateImmutable(cnf *ChainStateConfig) (*ChainState, error) {
	return newChainStateEx(cnf, false)
}

func newChainStateEx(cnf *ChainStateConfig, isSubBranchTemporary bool) (*ChainState, error) {

	sid := rand.Uint64()
	state := &ChainState{
		sid:              sid,
		config:           cnf,
		base:             nil,
		childs:           make([]*ChainState, 0),
		referBlockHeight: 0,
		referBlockHash:   nil,
		ldb:              nil,
		memdb:            nil,
		statusMux:        &sync.RWMutex{},
	}

	// 建立数据库 db
	if isSubBranchTemporary {
		useldb, e := leveldb.OpenFile(cnf.Datadir, nil)
		if e != nil {
			return nil, e
		}
		state.ldb = useldb
	} else {
		state.memdb = &sync.Map{}
	}

	// ok
	return state, nil
}

func (s ChainState) IsImmutable() bool {
	s.statusMux.RLock()
	defer s.statusMux.RUnlock()

	if s.ldb != nil {
		return true
	}
	return false
}

func (s *ChainState) ForkSubChild() (*ChainState, error) {
	s.statusMux.Lock()
	defer s.statusMux.Unlock()

	// create
	sub, e := newChainStateEx(s.config, true)
	if e != nil {
		return nil, e
	}
	// set base
	sub.base = s
	// add childs
	s.childs = append(s.childs, sub)
	// ok
	return sub, nil
}

func (s *ChainState) ForkNextBlock(hei uint64, hx fields.Hash) (*ChainState, error) {
	// create
	sub, e := s.ForkSubChild()
	if e != nil {
		return nil, e
	}

	s.statusMux.Lock()
	defer s.statusMux.Unlock()

	// set block
	sub.referBlockHeight = hei
	sub.referBlockHash = hx
	// ok
	return sub, nil
}

// 获取指向的区块
func (s ChainState) GetReferBlock() (uint64, fields.Hash) {
	s.statusMux.RLock()
	defer s.statusMux.RUnlock()

	return s.referBlockHeight, s.referBlockHash
}

// 获得父级状态
func (s ChainState) GetParent() *ChainState {
	s.statusMux.RLock()
	defer s.statusMux.RUnlock()

	return s.base
}

// 获得所有子状态
func (s ChainState) GetChilds() []*ChainState {
	s.statusMux.RLock()
	defer s.statusMux.RUnlock()

	return s.childs
}

// 销毁，包括删除所有子状态、缓存、状态数据等
func (s *ChainState) Destory() {
	s.statusMux.Lock()
	defer s.statusMux.Unlock()
	// clean
	s.memdb = nil
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
