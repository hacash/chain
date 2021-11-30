package chainstatev3

import (
	"bytes"
	"github.com/hacash/core/blocks"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfacev3"
	"github.com/hacash/core/stores"
)

type LatestStatus struct {
	// 不成熟的区块哈希
	ImmatureBlockHashCount fields.VarUint2
	ImmatureBlockHashs     []fields.Hash

	// 成熟的区块头
	ImmutableBlockHeadMetaIsExist fields.Bool
	ImmutableBlockHeadMeta        interfacev3.Block

	// 最新的区块头
	LatestBlockHashIsExist fields.Bool
	LatestBlockHash        fields.Hash

	// 最新区块钻石
	LatestDiamondIsExist fields.Bool
	LatestDiamond        *stores.DiamondSmelt
}

func (p *LatestStatus) Size() uint32 {
	size := p.ImmatureBlockHashCount.Size()
	size += uint32(p.ImmatureBlockHashCount) * fields.HashSize
	// ImmutableBlockHeadMetaIsExist
	size += p.ImmutableBlockHeadMetaIsExist.Size()
	if p.ImmutableBlockHeadMetaIsExist.Check() {
		size += blocks.BlockHeadSize + blocks.BlockMetaSizeV1
	}
	// LatestBlockHashIsExist
	size += p.LatestBlockHashIsExist.Size()
	if p.LatestBlockHashIsExist.Check() {
		size += fields.HashSize
	}
	// LatestDiamondIsExist
	size += p.LatestDiamondIsExist.Size()
	if p.LatestDiamondIsExist.Check() {
		size += p.LatestDiamond.Size()
	}
	return size
}

func (p *LatestStatus) Serialize() ([]byte, error) {
	var e error = nil
	var bt []byte = nil
	var buf = bytes.NewBuffer(nil)
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
	// ImmatureBlockHashCount
	seek, e = p.ImmatureBlockHashCount.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	for i := 0; i < int(p.ImmatureBlockHashCount); i++ {
		seek, e = p.ImmatureBlockHashs[i].Parse(buf, seek)
		if e != nil {
			return 0, e
		}
	}
	// ImmutableBlockHeadMetaIsExist
	seek, e = p.ImmutableBlockHeadMetaIsExist.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	if p.ImmutableBlockHeadMetaIsExist.Check() {
		seek, e = p.ImmutableBlockHeadMeta.ParseExcludeTransactions(buf, seek)
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
	// LatestDiamondIsExist
	seek, e = p.LatestDiamondIsExist.Parse(buf, seek)
	if e != nil {
		return 0, e
	}
	if p.LatestDiamondIsExist.Check() {
		seek, e = p.LatestDiamond.Parse(buf, seek)
		if e != nil {
			return 0, e
		}
	}
	return seek, nil
}

/////////////////////////////////////////////

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

func (p *LatestStatus) GetImmatureBlockHashList() []fields.Hash {
	if p.ImmatureBlockHashCount > 0 {
		return p.ImmatureBlockHashs
	}
	return []fields.Hash{}
}

func (p *LatestStatus) GetLatestBlockHash() fields.Hash {
	if p.LatestBlockHashIsExist.Check() {
		return p.LatestBlockHash
	}
	return nil
}

func (p *LatestStatus) GetImmutableBlockHeadMeta() interfacev3.Block {
	if p.ImmutableBlockHeadMetaIsExist.Check() {
		return p.ImmutableBlockHeadMeta
	}
	return nil
}

func (p *LatestStatus) SetImmutableBlockHeadMeta(blkHeadMeta interfacev3.Block) {
	p.ImmutableBlockHeadMetaIsExist.Set(true)
	p.ImmutableBlockHeadMeta = blkHeadMeta
}

func (p *LatestStatus) SetImmatureBlockHashList(hxs []fields.Hash) {
	p.ImmatureBlockHashCount = fields.VarUint2(len(hxs))
	p.ImmatureBlockHashs = hxs
}

/////////////////////////////////////////////

func NewInitialLatestStatus() *LatestStatus {

	return &LatestStatus{
		ImmatureBlockHashCount: fields.VarUint2(0),
		ImmatureBlockHashs:     nil,

		ImmutableBlockHeadMetaIsExist: fields.CreateBool(false),
		ImmutableBlockHeadMeta:        nil,

		LatestBlockHashIsExist: fields.CreateBool(false),
		LatestBlockHash:        nil,

		LatestDiamondIsExist: fields.CreateBool(false),
		LatestDiamond:        nil,
	}
}
