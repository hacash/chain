package blockstorev2

import (
	"encoding/binary"
	"fmt"
	"github.com/hacash/core/stores"
)

// btc move log

// Number of data pages, 100 pieces of data per page
func (cs *BlockStore) GetBTCMoveLogTotalPage() (int, error) {
	if cs.btcmovelogTotalPage >= 0 {
		return cs.btcmovelogTotalPage, nil
	}
	// Read from storage
	bts, _ := cs.btcmovelogDB.Get([]byte("total_page"))
	//fmt.Println("---------bts----", bts)
	if len(bts) != 4 {
		cs.btcmovelogTotalPage = 0 // No data
	} else {
		pg := binary.BigEndian.Uint32(bts)
		cs.btcmovelogTotalPage = int(pg)
	}
	// return
	return cs.btcmovelogTotalPage, nil
}

// Get data page
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

	// analysis
	return stores.SatoshiGenesisPageParse(dtbts, 0), nil
}

// Save data page
func (cs *BlockStore) SaveBTCMoveLogPageData(svpage int, list []*stores.SatoshiGenesis) error {
	// Save page number
	if svpage >= cs.btcmovelogTotalPage {
		cs.btcmovelogTotalPage = svpage
		pgk := []byte("total_page")
		pgd := make([]byte, 4)
		binary.BigEndian.PutUint32(pgd, uint32(svpage))
		//fmt.Println("-------------", svpage, pgd)
		cs.btcmovelogDB.Set(pgk, pgd)
	}
	// Save content
	datas := stores.SatoshiGenesisPageSerialize(list)
	//fmt.Println(strings.Join(stores.SatoshiGenesisPageSerializeForShow(list), " | "))
	//fmt.Println(stores.SatoshiGenesisPageParse(datas, 0))
	//fmt.Println("-------cs.btcmovelogDB.Set(key, datas)------", len(datas), stores.SatoshiGenesisPageParse(datas, 0), strings.Join(stores.SatoshiGenesisPageSerializeForShow(list), " | "))
	key := []byte(fmt.Sprintf("page_data_%d", svpage))
	return cs.btcmovelogDB.Set(key, datas)
}
