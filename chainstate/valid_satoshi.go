package chainstate

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
	"github.com/hacash/core/stores"
	"github.com/hacash/mint/coinbase"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
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
var btcMoveLocalLogsCachePageData []string = nil

func readSatoshiGenesisByLocalLogs(store interfaces.BlockStore, trsno int64) *stores.SatoshiGenesis {
	var limit = stores.SatoshiGenesisLogStorePageLimit // limit 100
	var logitemstr string = ""
	readpage := int((trsno-1)/int64(limit)) + 1
	offset := int((trsno - 1) % int64(limit))
	// 从缓存读取
	if readpage == btcMoveLocalLogsCachePage {
		logitemstr = btcMoveLocalLogsCachePageData[offset]
		return parseSatoshiGenesisByItemString(logitemstr, trsno)
	}
	// 从日志读取
	pagedata, err := store.GetBTCMoveLogPageData(readpage)
	if err != nil {
		return nil // 日志不存在
	}
	// 缓存
	if len(pagedata) == limit {
		btcMoveLocalLogsCachePage = readpage
		btcMoveLocalLogsCachePageData = pagedata
	}
	// 获取
	if offset >= len(pagedata) {
		return nil // 超出范围
	}
	logitemstr = pagedata[offset]
	// 解析
	return parseSatoshiGenesisByItemString(logitemstr, trsno)
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
	return parseSatoshiGenesisByItemString(logitemstr, trsno)
}

// 解析日志
func parseSatoshiGenesisByItemString(logitemstr string, trsno int64) *stores.SatoshiGenesis {
	// 开始解析
	logitemstr = strings.Replace(logitemstr, " ", "", -1)
	dts := strings.Split(logitemstr, ",")
	if len(dts) != 8 {
		return nil
	}
	nums := make([]int64, 6)
	for i := 0; i < 6; i++ {
		n, e := strconv.ParseInt(dts[i], 10, 0)
		if e != nil {
			return nil
		}
		nums[i] = n
	}
	if nums[0] != trsno {
		return nil // 标号对不上
	}
	// 检查地址 和 txhx
	addr, ae := fields.CheckReadableAddress(dts[6])
	if ae != nil {
		return nil
	}
	trshx, te := hex.DecodeString(dts[7])
	if te != nil {
		return nil
	}
	if len(trshx) != 32 {
		return nil
	}
	// 检查转账数量
	ttb := nums[3]
	btcs := nums[4]
	if btcs < 1 && btcs > 10000 {
		return nil // 转移的比特币最小一枚，最大 10000 枚（超过10000的按1000计算）
	}
	var totalAddHAC int64 = 0
	for i := ttb + 1; i <= ttb+btcs; i++ {
		totalAddHAC += coinbase.MoveBtcCoinRewardNumber(i)
	}
	if totalAddHAC != nums[5] {
		return nil // 增发的hac对不上
	}
	// 生成
	genesis := stores.SatoshiGenesis{
		TransferNo:               fields.VarUint4(nums[0]),     // 转账流水编号
		BitcoinBlockHeight:       fields.VarUint4(nums[1]),     // 转账的比特币区块高度
		BitcoinBlockTimestamp:    fields.VarUint4(nums[2]),     // 转账的比特币区块时间戳
		BitcoinEffectiveGenesis:  fields.VarUint4(nums[3]),     // 在这笔之前已经成功转移的比特币数量
		BitcoinQuantity:          fields.VarUint4(nums[4]),     // 本笔转账的比特币数量（单位：枚）
		AdditionalTotalHacAmount: fields.VarUint4(totalAddHAC), // 本次转账[总共]应该增发的 hac 数量 （单位：枚）
		OriginAddress:            *addr,
		BitcoinTransferHash:      trshx,
	}

	// 返回
	return &genesis
}
