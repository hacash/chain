package chainstatev3

import (
	"github.com/hacash/core/fields"
	"github.com/hacash/core/stores"
)

//
func (cs *ChainState) Channel(channelid fields.ChannelId) (*stores.Channel, error) {
	value, ok, e := cs.find(KeySuffixType_channel, channelid)
	if e != nil {
		return nil, e
	}
	if !ok {
		return nil, nil // not find
	}
	// parse
	var stoitem stores.Channel
	_, e = stoitem.Parse(value, 0)
	if e != nil {
		return nil, e // error
	}
	// return ok
	return &stoitem, nil
}

//
func (cs *ChainState) ChannelCreate(channel_id fields.ChannelId, channel *stores.Channel) error {
	stodatas, e := channel.Serialize()
	if e != nil {
		return e // error
	}
	// do save
	return cs.save(KeySuffixType_channel, channel_id, stodatas)
}

//
func (cs *ChainState) ChannelUpdate(channel_id fields.ChannelId, channel *stores.Channel) error {
	return cs.ChannelCreate(channel_id, channel)
}

//
func (cs *ChainState) ChannelDelete(channel_id fields.ChannelId) error {
	return cs.delete(KeySuffixType_btclend, channel_id)
}
