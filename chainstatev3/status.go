package chainstatev3

import (
	"fmt"
	"github.com/hacash/core/interfaces"
)

const (
	status_key_latest    = "latest"
	status_key_immutable = "immutable"
)

func (cs *ChainState) StatusRead(name string) ([]byte, error) {
	value, ok, e := cs.find(KeySuffixType_statuskv, []byte(name))
	if e != nil {
		return nil, e
	}
	if !ok {
		return nil, nil // not find
	}
	return value, nil
}

func (cs *ChainState) StatusSet(name string, value []byte) error {
	return cs.save(KeySuffixType_statuskv, []byte(name), value)
}

/////////////////////////////////////////

func (cs *ChainState) LatestStatusRead() (interfaces.LatestStatus, error) {
	value, e := cs.StatusRead(status_key_latest)
	if e != nil {
		return nil, e
	}
	if value == nil {
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

func (cs *ChainState) LatestStatusSet(status interfaces.LatestStatus) error {
	// 保存
	datas, e := status.Serialize()
	if e != nil {
		return e
	}
	// do save
	return cs.StatusSet(status_key_latest, datas)
}

func (cs *ChainState) ImmutableStatusSet(status interfaces.ImmutableStatus) error {
	// 保存
	datas, e := status.Serialize()
	if e != nil {
		return e
	}
	// do save
	return cs.StatusSet(status_key_immutable, datas)
}

func (cs *ChainState) ImmutableStatusRead() (interfaces.ImmutableStatus, error) {
	value, e := cs.StatusRead(status_key_immutable)
	if e != nil {
		return nil, e
	}
	if value == nil {
		// not find 返回初始状态
		return NewInitialImmutableStatus(), nil
	}

	var stoitem ImmutableStatus
	_, e = stoitem.Parse(value, 0)
	if e != nil {
		return nil, e // error
	}
	// return ok
	return &stoitem, nil
}

//////////////////////////////////////////////////////

func (cs *ChainState) GetPending() interfaces.PendingStatus {
	if cs.pending == nil {
		panic(fmt.Errorf("pending cannot be nil."))
	}
	return cs.pending

}

func (cs *ChainState) SetPending(pd interfaces.PendingStatus) error {
	cs.pending = pd
	return nil
}
