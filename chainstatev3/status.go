package chainstatev3

import (
	"fmt"
	"github.com/hacash/core/interfacev3"
)

func (cs *ChainState) LatestStatusRead() (interfacev3.LatestStatus, error) {
	if cs.lastStatusCache != nil {
		return cs.lastStatusCache, nil
	}
	value, ok, e := cs.find(KeySuffixType_laststatus, []byte{1})
	if e != nil {
		return nil, e
	}
	if !ok {
		// not find 返回初始状态
		return NewInitialLatestStatus(), nil
	}

	var stoitem LatestStatus
	_, e = stoitem.Parse(value, 0)
	if e != nil {
		return nil, e // error
	}
	// return ok
	return &stoitem, nil
}

func (cs *ChainState) LatestStatusSet(status interfacev3.LatestStatus) error {
	// 缓存状态
	cs.lastStatusCache = status
	// 保存
	datas, e := status.Serialize()
	if e != nil {
		return e
	}
	// do save
	return cs.save(KeySuffixType_laststatus, []byte{1}, datas)
}

func (cs *ChainState) GetPending() interfacev3.PendingStatus {
	if cs.pending == nil {
		panic(fmt.Errorf("pending cannot be nil."))
	}
	return cs.pending

}

func (cs *ChainState) SetPending(pd interfacev3.PendingStatus) error {
	cs.pending = pd
	return nil
}
