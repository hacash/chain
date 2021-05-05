package blockstore

import (
	"encoding/binary"
	"fmt"
	"github.com/hacash/core/stores"
)

// btc move log

// 数据页数，每页100条数据
func (cs *BlockStore) GetBTCMoveLogTotalPage() (int, error) {
	if cs.btcmovelogTotalPage >= 0 {
		return cs.btcmovelogTotalPage, nil
	}
	// 从储存读取
	bts, _ := cs.btcmovelogDB.Get([]byte("total_page"))
	//fmt.Println("---------bts----", bts)
	if len(bts) != 4 {
		cs.btcmovelogTotalPage = 0 // 无数据
	} else {
		pg := binary.BigEndian.Uint32(bts)
		cs.btcmovelogTotalPage = int(pg)
	}
	// 返回
	return cs.btcmovelogTotalPage, nil
}

// 获取数据页
func (cs *BlockStore) GetBTCMoveLogPageData(page int) ([]*stores.SatoshiGenesis, error) {
	realpage, e0 := cs.GetBTCMoveLogTotalPage()
	if e0 != nil {
		return nil, e0
	}
	if page > realpage {
		return nil, fmt.Errorf("overflow data page")
	}
	dtbts, e2 := cs.btcmovelogDB.Get([]byte(fmt.Sprintf("page_data_%d", page)))
	if e2 != nil {
		return []*stores.SatoshiGenesis{}, nil
	}
	// 解析
	return stores.SatoshiGenesisPageParse(dtbts, 0), nil
}

// 保存数据页
func (cs *BlockStore) SaveBTCMoveLogPageData(svpage int, list []*stores.SatoshiGenesis) error {
	// 保存页码
	if svpage >= cs.btcmovelogTotalPage {
		cs.btcmovelogTotalPage = svpage
		pgk := []byte("total_page")
		pgd := make([]byte, 4)
		binary.BigEndian.PutUint32(pgd, uint32(svpage))
		//fmt.Println("-------------", svpage, pgd)
		cs.btcmovelogDB.Set(pgk, pgd)
	}
	// 保存内容
	datas := stores.SatoshiGenesisPageSerialize(list)
	//fmt.Println(strings.Join(stores.SatoshiGenesisPageSerializeForShow(list), " | "))
	//fmt.Println(stores.SatoshiGenesisPageParse(datas, 0))
	//fmt.Println("-------cs.btcmovelogDB.Set(key, datas)------", len(datas), stores.SatoshiGenesisPageParse(datas, 0), strings.Join(stores.SatoshiGenesisPageSerializeForShow(list), " | "))
	key := []byte(fmt.Sprintf("page_data_%d", svpage))
	return cs.btcmovelogDB.Set(key, datas)
}
