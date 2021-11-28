package chainstatev3

import (
	"github.com/hacash/core/fields"
	"github.com/hacash/core/stores"
)

func (cs *ChainState) Lockbls(lkid fields.LockblsId) (*stores.Lockbls, error) {
	value, ok, e := cs.find(KeySuffixType_lockbls, lkid)
	if e != nil {
		return nil, e
	}
	if !ok {
		return nil, nil // not find
	}
	// parse
	var stoitem stores.Lockbls
	_, e = stoitem.Parse(value, 0)
	if e != nil {
		return nil, e // error
	}
	// return ok
	return &stoitem, nil
}

func (cs *ChainState) LockblsCreate(lkid fields.LockblsId, stoitem *stores.Lockbls) error {
	stodatas, e := stoitem.Serialize()
	if e != nil {
		return e // error
	}
	// do save
	return cs.save(KeySuffixType_lockbls, lkid, stodatas)
}

func (cs *ChainState) LockblsUpdate(lkid fields.LockblsId, stoitem *stores.Lockbls) error {
	return cs.LockblsCreate(lkid, stoitem)
}

func (cs *ChainState) LockblsDelete(lkid fields.LockblsId) error {
	return cs.delete(KeySuffixType_lockbls, lkid)
}
