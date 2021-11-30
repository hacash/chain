package chainstatev3

import (
	"bytes"
	"github.com/hacash/core/blocks"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/stores"
)

type PendingStatus struct {
	blockHeight fields.BlockHeight

	blockHashIsExist fields.Bool
	blockHash        fields.Hash

	blockHeadMetaIsExist fields.Bool
	blockHeadMeta        interfaces.BlockHeadMetaRead

	waitingSubmitDiamondIsExist fields.Bool
	waitingSubmitDiamond        *stores.DiamondSmelt
}

func (p *PendingStatus) Size() uint32 {
	size := fields.BlockHeightSize
	size += p.blockHashIsExist.Size()
	if p.blockHashIsExist.Check() {
		size += fields.HashSize
	}
	size += p.blockHeadMetaIsExist.Size()
	if p.blockHeadMetaIsExist.Check() {
		size += blocks.BlockHeadSize + blocks.BlockMetaSizeV1
	}
	size += p.waitingSubmitDiamondIsExist.Size()
	if p.waitingSubmitDiamondIsExist.Check() {
		size += p.waitingSubmitDiamond.Size()
	}
	return size
}

func (p *PendingStatus) Serialize() ([]byte, error) {
	var e error = nil
	var bt []byte = nil
	var buf = bytes.NewBuffer(nil)
	// BlockHeight
	bt, e = p.blockHeight.Serialize()
	if e != nil {
		return nil, e
	}
	buf.Write(bt)
	// blockHashIsExist
	bt, e = p.blockHashIsExist.Serialize()
	if e != nil {
		return nil, e
	}
	buf.Write(bt)
	if p.blockHashIsExist.Check() {
		bt, e = p.blockHash.Serialize()
		if e != nil {
			return nil, e
		}
		buf.Write(bt)
	}
	// blockHeadMetaIsExist
	bt, e = p.blockHeadMetaIsExist.Serialize()
	if e != nil {
		return nil, e
	}
	buf.Write(bt)
	if p.blockHeadMetaIsExist.Check() {
		bt, e = p.blockHeadMeta.SerializeExcludeTransactions()
		if e != nil {
			return nil, e
		}
		buf.Write(bt)
	}
	// waitingSubmitDiamondIsExist
	bt, e = p.waitingSubmitDiamondIsExist.Serialize()
	if e != nil {
		return nil, e
	}
	buf.Write(bt)
	if p.waitingSubmitDiamondIsExist.Check() {
		bt, e = p.waitingSubmitDiamond.Serialize()
		if e != nil {
			return nil, e
		}
		buf.Write(bt)
	}
	// ok
	return buf.Bytes(), nil
}

func (p *PendingStatus) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error = nil
	seek, e = p.blockHeight.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	// blockHashIsExist
	seek, e = p.blockHashIsExist.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	if p.blockHashIsExist.Check() {
		seek, e = p.blockHash.Parse(buf, seek)
		if e != nil {
			return 0, e
		}
	}
	// blockHeadMetaIsExist
	seek, e = p.blockHeadMetaIsExist.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	if p.blockHeadMetaIsExist.Check() {
		p.blockHeadMeta, seek, e = blocks.ParseExcludeTransactions(buf, seek)
		if e != nil {
			return 0, e
		}
	}
	// waitingSubmitDiamondIsExist
	seek, e = p.waitingSubmitDiamondIsExist.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	if p.waitingSubmitDiamondIsExist.Check() {
		seek, e = p.waitingSubmitDiamond.Parse(buf, seek)
		if e != nil {
			return 0, e
		}
	}
	return seek, nil
}

/////////////////////////////////////////////

func (p *PendingStatus) GetPendingBlockHead() interfaces.BlockHeadMetaRead {
	if p.blockHeadMetaIsExist.Check() {
		return p.blockHeadMeta
	}
	return nil
}

func (p *PendingStatus) GetPendingBlockHeight() uint64 {
	if p.blockHeadMetaIsExist.Check() {
		return p.blockHeadMeta.GetHeight()
	}
	return uint64(p.blockHeight)
}

func (p *PendingStatus) GetPendingBlockHash() fields.Hash {
	if p.blockHeadMetaIsExist.Check() {
		return p.blockHeadMeta.Hash()
	}
	return p.blockHash
}

func (p *PendingStatus) GetWaitingSubmitDiamond() *stores.DiamondSmelt {
	if p.waitingSubmitDiamondIsExist.Check() {
		return p.waitingSubmitDiamond
	}
	return nil
}

func (p *PendingStatus) SetWaitingSubmitDiamond(diamond *stores.DiamondSmelt) {
	p.waitingSubmitDiamondIsExist.Set(true)
	p.waitingSubmitDiamond = diamond
}

func (p *PendingStatus) ClearWaitingSubmitDiamond() {
	p.waitingSubmitDiamondIsExist.Set(false)
	p.waitingSubmitDiamond = nil
}

/////////////////////////////////////////////

func NewPendingStatus(hei uint64, hx fields.Hash, blockhead interfaces.BlockHeadMetaRead) *PendingStatus {
	ins := &PendingStatus{
		blockHeight:                 0,
		blockHashIsExist:            fields.CreateBool(false),
		blockHash:                   nil,
		blockHeadMetaIsExist:        fields.CreateBool(false),
		blockHeadMeta:               nil,
		waitingSubmitDiamondIsExist: fields.CreateBool(false),
		waitingSubmitDiamond:        nil,
	}
	if blockhead != nil {
		ins.blockHeadMetaIsExist = fields.CreateBool(true)
		ins.blockHeadMeta = blockhead
	}
	if hei > 0 {
		ins.blockHeight = fields.BlockHeight(hei)
	}
	if hx != nil {
		ins.blockHashIsExist = fields.CreateBool(true)
		ins.blockHash = hx
	}
	return ins
}
