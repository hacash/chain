package blockstorev3

import (
	"encoding/binary"
	"fmt"
	"github.com/hacash/chain/leveldb"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/stores"
	"github.com/hacash/mint/coinbase"
	"io/ioutil"
	"net/http"
	"time"
)

// btc move log

// 数据页数，每页100条数据
func (cs *BlockStore) GetBTCMoveLogTotalPage() (int, error) {
	cs.statusMux.RLock()
	var page = cs.btcmovelogTotalPage
	cs.statusMux.RUnlock()
	if page >= 0 {
		return page, nil
	}

	ldb, e := cs.getDB()
	if e != nil {
		return -1, e
	}

	cs.statusMux.Lock()
	defer cs.statusMux.Unlock()

	// 从储存读取
	bts, e := ldb.Get([]byte("btc_move_log_total_page"), nil)
	if e != nil {
		if e == leveldb.ErrNotFound {
			cs.btcmovelogTotalPage = 0 // 无数据
			return 0, nil
		}
		return -1, e
	}
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

	ldb, e := cs.getDB()
	if e != nil {
		return nil, e
	}

	// read
	pgkey := []byte(fmt.Sprintf("btc_move_log_page_data_%d", page))
	dtbts, e := ldb.Get(pgkey, nil)
	if e != nil {
		if e == leveldb.ErrNotFound {
			return []*stores.SatoshiGenesis{}, nil
		}
		return nil, e
	}
	// 解析
	return stores.SatoshiGenesisPageParse(dtbts, 0), nil
}

// 保存数据页
func (cs *BlockStore) SaveBTCMoveLogPageData(svpage int, list []*stores.SatoshiGenesis) error {

	ldb, e := cs.getDB()
	if e != nil {
		return e
	}

	cs.statusMux.Lock()
	defer cs.statusMux.Unlock()

	// 保存页码
	if svpage >= cs.btcmovelogTotalPage {
		cs.btcmovelogTotalPage = svpage
		pgk := []byte("btc_move_log_total_page")
		pgd := make([]byte, 4)
		binary.BigEndian.PutUint32(pgd, uint32(svpage))
		//fmt.Println("-------------", svpage, pgd)
		e := ldb.Put(pgk, pgd, nil)
		if e != nil {
			return e
		}
	}
	// 保存内容
	datas := stores.SatoshiGenesisPageSerialize(list)
	pgkey := []byte(fmt.Sprintf("btc_move_log_page_data_%d", svpage))
	return ldb.Put(pgkey, datas, nil)
}

///////////////////////////////////////////

// 获取已验证的BTC转移日志，返回获取的内容以及是否需要验证
func (cs *BlockStore) LoadValidatedSatoshiGenesis(trsno int64) (*stores.SatoshiGenesis, bool) {
	//fmt.Println("LoadValidatedSatoshiGenesis: trsno:", trsno)

	var mustcheck = false
	var genesis *stores.SatoshiGenesis = nil
	//fmt.Println(cs.config.BTCMoveCheckEnable, cs.config.BTCMoveCheckLogsURL)
	if cs.config.BTCMoveCheckEnable {
		mustcheck = true
		// 先从日志读取转移记录
		genesis = readSatoshiGenesisByLocalLogs(cs, trsno)
		if genesis == nil {
			// 日志里没有，再从 URL 读取
			if cs.config.BTCMoveCheckEnable {
				genesis = readSatoshiGenesisOneByUrl(cs.config.DownloadBTCMoveLogUrl, trsno)
			}
		} else {
			fmt.Printf("[Satoshi genesis] load from local database, Trsno: %d, BTC: %d, Address: %s.\n", genesis.TransferNo, genesis.BitcoinQuantity, genesis.OriginAddress.ToReadable())
		}
	}
	// 返回
	return genesis, mustcheck
}

// 读取缓存
var btcMoveLocalLogsCachePage int = -1
var btcMoveLocalLogsCachePageData []*stores.SatoshiGenesis = nil

func readSatoshiGenesisByLocalLogs(store interfaces.BlockStore, trsno int64) *stores.SatoshiGenesis {
	var limit = stores.SatoshiGenesisLogStorePageLimit // limit 100
	readpage := int((trsno-1)/int64(limit)) + 1
	offset := int((trsno - 1) % int64(limit))
	// 从缓存读取
	if readpage == btcMoveLocalLogsCachePage {
		return btcMoveLocalLogsCachePageData[offset]
	}
	// 从日志读取
	pagedata, err := store.GetBTCMoveLogPageData(readpage)
	if err != nil {
		return nil // 日志不存在
	}
	// 获取
	if offset >= len(pagedata) {
		return nil // 超出范围
	}
	// 解析
	if len(pagedata) == limit {
		btcMoveLocalLogsCachePage = readpage
		btcMoveLocalLogsCachePageData = pagedata
	}
	// ok
	return pagedata[offset]
}

func readSatoshiGenesisOneByUrl(url string, trsno int64) *stores.SatoshiGenesis {
	if len(url) == 0 {
		return nil // error
	}
	client := http.Client{
		Timeout: time.Duration(10 * time.Second),
	}
	url += fmt.Sprintf("?trsno=%d", trsno)
	fmt.Println("[Satoshi genesis] load check url:", url)
	resp, err := client.Get(url)
	if err != nil {
		//fmt.Println("read Validated SatoshiGenesisByUrl return error:", err.Error())
		return nil
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	logitemstr := string(body)
	fmt.Println("[Satoshi genesis] got data by url:", logitemstr)

	// 解析
	return coinbase.ParseSatoshiGenesisByItemString(logitemstr, trsno)
}
