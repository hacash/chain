package chainstatev2

import (
	"encoding/binary"
	"fmt"
	"github.com/hacash/core/interfacev2"
	"github.com/hacash/core/stores"
	"github.com/hacash/mint/coinbase"
	"io/ioutil"
	"net/http"
	"time"
)

func (cs *ChainState) SetInMemTxPool(stat bool) {
	cs.isInTxPool = stat
}

func (cs *ChainState) IsInMemTxPool() bool {
	return cs.isInTxPool
}

// block data store
func (cs *ChainState) SaveMoveBTCBelongTxHash(trsno uint32, txhash []byte) error {
	// save trsno name by number key
	numberkey := make([]byte, 4)
	binary.BigEndian.PutUint32(numberkey, trsno)
	query, e4 := cs.movebtcDB.CreateNewQueryInstance(numberkey)
	if e4 != nil {
		return e4
	}
	defer query.Destroy()
	e5 := query.Save(txhash)
	if e5 != nil {
		return e5
	}
	return nil
}

// block data store
func (cs *ChainState) ReadMoveBTCTxHashByTrsNo(trsno uint32) ([]byte, error) {
	return cs.ReadMoveBTCTxHashByNumber(trsno)
}
func (cs *ChainState) ReadMoveBTCTxHashByNumber(number uint32) ([]byte, error) {
	// find by number key
	numberkey := make([]byte, 4)
	binary.BigEndian.PutUint32(numberkey, number)
	query1, e1 := cs.movebtcDB.CreateNewQueryInstance(numberkey)
	if e1 != nil {
		return nil, e1
	}
	defer query1.Destroy()
	belongtxhash, e3 := query1.Find()
	if e3 != nil {
		return nil, e3
	}
	if belongtxhash == nil && cs.base != nil {
		// 向上查询
		return cs.base.ReadMoveBTCTxHashByNumber(number)
	}

	if len(belongtxhash) > 0 && len(belongtxhash) != 32 {
		return nil, fmt.Errorf("move btc store file break.")
	}
	return belongtxhash, nil

}

// 获取已验证的BTC转移日志，返回获取的内容以及是否需要验证
func (cs *ChainState) LoadValidatedSatoshiGenesis(trsno int64) (*stores.SatoshiGenesis, bool) {
	//fmt.Println("LoadValidatedSatoshiGenesis: trsno:", trsno)

	var mustcheck = false
	var genesis *stores.SatoshiGenesis = nil
	if cs.isInTxPool {
		mustcheck = true
	}
	//fmt.Println(cs.config.BTCMoveCheckEnable, cs.config.BTCMoveCheckLogsURL)
	if cs.config.BTCMoveCheckEnable {
		mustcheck = true
		// 先从日志读取转移记录
		genesis = readSatoshiGenesisByLocalLogs(cs.BlockStore(), trsno)
		if genesis == nil {
			// 日志里没有，再从 URL 读取
			if len(cs.config.BTCMoveCheckLogsURL) > 0 {
				genesis = readSatoshiGenesisByUrl(cs.config.BTCMoveCheckLogsURL, trsno)
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

func readSatoshiGenesisByLocalLogs(store interfacev2.BlockStore, trsno int64) *stores.SatoshiGenesis {
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

func readSatoshiGenesisByUrl(url string, trsno int64) *stores.SatoshiGenesis {
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
