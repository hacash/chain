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

// Number of data pages, 100 pieces of data per page
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

	// Read from storage
	bts, e := ldb.Get([]byte("btc_move_log_total_page"), nil)
	if e != nil {
		if e == leveldb.ErrNotFound {
			cs.btcmovelogTotalPage = 0 // No data
			return 0, nil
		}
		return -1, e
	}
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
	// analysis
	return stores.SatoshiGenesisPageParse(dtbts, 0), nil
}

// Save data page
func (cs *BlockStore) SaveBTCMoveLogPageData(svpage int, list []*stores.SatoshiGenesis) error {

	ldb, e := cs.getDB()
	if e != nil {
		return e
	}

	cs.statusMux.Lock()
	defer cs.statusMux.Unlock()

	// Save page number
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
	// Save content
	datas := stores.SatoshiGenesisPageSerialize(list)
	pgkey := []byte(fmt.Sprintf("btc_move_log_page_data_%d", svpage))
	return ldb.Put(pgkey, datas, nil)
}

///////////////////////////////////////////

// Get the verified BTC transfer log, return the obtained content and whether verification is required
func (cs *BlockStore) LoadValidatedSatoshiGenesis(trsno int64) (*stores.SatoshiGenesis, bool) {
	//fmt.Println("LoadValidatedSatoshiGenesis: trsno:", trsno)

	var mustcheck = false
	var genesis *stores.SatoshiGenesis = nil
	//fmt.Println(cs.config.BTCMoveCheckEnable, cs.config.BTCMoveCheckLogsURL)
	if cs.config.BTCMoveCheckEnable {
		mustcheck = true
		// Read the transfer record from the log first
		genesis = readSatoshiGenesisByLocalLogs(cs, trsno)
		if genesis == nil {
			// No in the log, read from the URL again
			if cs.config.BTCMoveCheckEnable {
				genesis = readSatoshiGenesisOneByUrl(cs.config.DownloadBTCMoveLogUrl, trsno)
			}
		} else {
			fmt.Printf("[Satoshi genesis] load from local database, Trsno: %d, BTC: %d, Address: %s.\n", genesis.TransferNo, genesis.BitcoinQuantity, genesis.OriginAddress.ToReadable())
		}
	}
	// return
	return genesis, mustcheck
}

// Read cache
var btcMoveLocalLogsCachePage int = -1
var btcMoveLocalLogsCachePageData []*stores.SatoshiGenesis = nil

func readSatoshiGenesisByLocalLogs(store interfaces.BlockStore, trsno int64) *stores.SatoshiGenesis {
	var limit = stores.SatoshiGenesisLogStorePageLimit // limit 100
	readpage := int((trsno-1)/int64(limit)) + 1
	offset := int((trsno - 1) % int64(limit))
	// Read from cache
	if readpage == btcMoveLocalLogsCachePage {
		return btcMoveLocalLogsCachePageData[offset]
	}
	// Read from log
	pagedata, err := store.GetBTCMoveLogPageData(readpage)
	if err != nil {
		return nil // log does not exist
	}
	// obtain
	if offset >= len(pagedata) {
		return nil // Out of range
	}
	// analysis
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

	// analysis
	return coinbase.ParseSatoshiGenesisByItemString(logitemstr, trsno)
}
