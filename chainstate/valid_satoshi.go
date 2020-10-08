package chainstate

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/hacash/core/fields"
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
	_, e5 := query.Save(txhash)
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
		// 递归查询
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
		if len(cs.config.BTCMoveCheckLogsURL) > 0 {
			genesis = readSatoshiGenesisByUrl(cs.config.BTCMoveCheckLogsURL, trsno)
		}
	}
	// 返回
	return genesis, mustcheck
}

func readSatoshiGenesisByUrl(url string, trsno int64) *stores.SatoshiGenesis {
	client := http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	url += fmt.Sprintf("?trsno=%d", trsno)
	fmt.Println("[satoshi genesis] load check:", url)
	resp, err := client.Get(url)
	if err != nil {
		//fmt.Println("read Validated SatoshiGenesisByUrl return error:", err.Error())
		return nil
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("[satoshi genesis] got data:", string(body))
	dts := strings.Split(string(body), ",")
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
	if btcs < 1 && btcs > 1000 {
		return nil // 转移的比特币最小一枚，最大 1000 枚（超过1000的按1000计算）
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
