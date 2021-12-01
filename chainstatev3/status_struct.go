package chainstatev3

import (
	"bytes"
	"github.com/hacash/core/blocks"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/genesis"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/stores"
)

type ImmutableStatus struct {

	// 成熟的区块头
	ImmutableBlockHeadMetaIsExist fields.Bool
	ImmutableBlockHeadMeta        interfaces.BlockHeadMetaRead

	// 最新的区块头
	LatestBlockHashIsExist fields.Bool
	LatestBlockHash        fields.Hash

	// 不成熟的区块哈希
	ImmatureBlockHashCount fields.VarUint2
	ImmatureBlockHashs     []fields.Hash
}

func (p *ImmutableStatus) Size() uint32 {
	// ImmutableBlockHeadMetaIsExist
	size := p.ImmutableBlockHeadMetaIsExist.Size()
	if p.ImmutableBlockHeadMetaIsExist.Check() {
		size += blocks.BlockHeadSize + blocks.BlockMetaSizeV1
	}
	// LatestBlockHashIsExist
	size += p.LatestBlockHashIsExist.Size()
	if p.LatestBlockHashIsExist.Check() {
		size += fields.HashSize
	}
	// ImmatureBlockHashCount
	size += p.ImmatureBlockHashCount.Size()
	size += uint32(p.ImmatureBlockHashCount) * fields.HashSize
	return size
}

func (p *ImmutableStatus) Serialize() ([]byte, error) {
	var e error = nil
	var bt []byte = nil
	var buf = bytes.NewBuffer(nil)
	// ImmutableBlockHeadMetaIsExist
	bt, e = p.ImmutableBlockHeadMetaIsExist.Serialize()
	if e != nil {
		return nil, e
	}
	buf.Write(bt)
	if p.ImmutableBlockHeadMetaIsExist.Check() {
		bt, e = p.ImmutableBlockHeadMeta.SerializeExcludeTransactions()
		if e != nil {
			return nil, e
		}
		buf.Write(bt)
	}
	// LatestBlockHashIsExist
	bt, e = p.LatestBlockHashIsExist.Serialize()
	if e != nil {
		return nil, e
	}
	buf.Write(bt)
	if p.LatestBlockHashIsExist.Check() {
		bt, e = p.LatestBlockHash.Serialize()
		if e != nil {
			return nil, e
		}
		buf.Write(bt)
	}
	// ImmatureBlockHashCount
	bt, e = p.ImmatureBlockHashCount.Serialize()
	if e != nil {
		return nil, e
	}
	buf.Write(bt)
	for i := 0; i < int(p.ImmatureBlockHashCount); i++ {
		bt, e = p.ImmatureBlockHashs[i].Serialize()
		if e != nil {
			return nil, e
		}
		buf.Write(bt)
	}
	// ok
	return buf.Bytes(), nil
}

func (p *ImmutableStatus) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error = nil
	// ImmutableBlockHeadMetaIsExist
	seek, e = p.ImmutableBlockHeadMetaIsExist.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	if p.ImmutableBlockHeadMetaIsExist.Check() {
		p.ImmutableBlockHeadMeta, seek, e = blocks.ParseExcludeTransactions(buf, seek)
		if e != nil {
			return 0, e
		}
	}
	// LatestBlockHashIsExist
	seek, e = p.LatestBlockHashIsExist.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	if p.LatestBlockHashIsExist.Check() {
		seek, e = p.LatestBlockHash.Parse(buf, seek)
		if e != nil {
			return 0, e
		}
	}
	// ImmatureBlockHashCount
	seek, e = p.ImmatureBlockHashCount.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	hxlen := int(p.ImmatureBlockHashCount)
	p.ImmatureBlockHashs = make([]fields.Hash, hxlen)
	for i := 0; i < hxlen; i++ {
		seek, e = p.ImmatureBlockHashs[i].Parse(buf, seek)
		if e != nil {
			return 0, e
		}
	}
	return seek, nil
}

func NewInitialImmutableStatus() *ImmutableStatus {
	return &ImmutableStatus{
		ImmutableBlockHeadMetaIsExist: fields.CreateBool(false),
		ImmutableBlockHeadMeta:        nil,
		LatestBlockHashIsExist:        fields.CreateBool(false),
		LatestBlockHash:               nil,
		ImmatureBlockHashCount:        0,
		ImmatureBlockHashs:            nil,
	}
}

func (p *ImmutableStatus) GetImmatureBlockHashList() []fields.Hash {
	if p.ImmatureBlockHashCount > 0 {
		return p.ImmatureBlockHashs
	}
	return []fields.Hash{}
}

func (p *ImmutableStatus) SetImmatureBlockHashList(list []fields.Hash) {
	p.ImmatureBlockHashCount = fields.VarUint2(len(list))
	p.ImmatureBlockHashs = list
}

func (p *ImmutableStatus) GetLatestBlockHash() fields.Hash {
	if p.LatestBlockHashIsExist.Check() {
		return p.LatestBlockHash
	}
	return nil
}

func (p *ImmutableStatus) SetLatestBlockHash(hx fields.Hash) {
	p.LatestBlockHashIsExist.Set(true)
	p.LatestBlockHash = hx
}

func (p *ImmutableStatus) GetImmutableBlockHeadMeta() interfaces.BlockHeadMetaRead {
	if p.ImmutableBlockHeadMetaIsExist.Check() {
		return p.ImmutableBlockHeadMeta
	}
	// 首次返回创始区块
	return genesis.GetGenesisBlock()
}

func (p *ImmutableStatus) SetImmutableBlockHeadMeta(head interfaces.BlockHeadMetaRead) {
	p.ImmutableBlockHeadMetaIsExist.Set(true)
	p.ImmutableBlockHeadMeta = head
}

//////////////////////////////////////////////////////////////////

type LatestStatus struct {
	// 最新区块钻石
	LatestDiamondIsExist fields.Bool
	LatestDiamond        *stores.DiamondSmelt
}

func (p *LatestStatus) Size() uint32 {
	// LatestDiamondIsExist
	size := p.LatestDiamondIsExist.Size()
	if p.LatestDiamondIsExist.Check() {
		size += p.LatestDiamond.Size()
	}
	return size
}

func (p *LatestStatus) Serialize() ([]byte, error) {
	var e error = nil
	var bt []byte = nil
	var buf = bytes.NewBuffer(nil)
	// LatestDiamondIsExist
	bt, e = p.LatestDiamondIsExist.Serialize()
	if e != nil {
		return nil, e
	}
	buf.Write(bt)
	if p.LatestDiamondIsExist.Check() {
		bt, e = p.LatestDiamond.Serialize()
		if e != nil {
			return nil, e
		}
		buf.Write(bt)
	}
	// ok
	return buf.Bytes(), nil
}

func (p *LatestStatus) Parse(buf []byte, seek uint32) (uint32, error) {
	var e error = nil
	// LatestDiamondIsExist
	seek, e = p.LatestDiamondIsExist.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	if p.LatestDiamondIsExist.Check() {
		p.LatestDiamond = &stores.DiamondSmelt{}
		seek, e = p.LatestDiamond.Parse(buf, seek)
		if e != nil {
			return 0, e
		}
	}
	return seek, nil
}

func (p *LatestStatus) SetLastestDiamond(diamond *stores.DiamondSmelt) {
	p.LatestDiamondIsExist.Set(true)
	p.LatestDiamond = diamond
}

func (p *LatestStatus) ReadLastestDiamond() *stores.DiamondSmelt {
	if p.LatestDiamondIsExist.Check() {
		return p.LatestDiamond
	}
	return nil
}

func NewInitialLatestStatus() *LatestStatus {
	return &LatestStatus{
		LatestDiamondIsExist: fields.CreateBool(false),
		LatestDiamond:        nil,
	}
}
