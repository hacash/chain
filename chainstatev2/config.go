package chainstatev2

import (
	"github.com/hacash/core/sys"
	"path"
)

type ChainStateConfig struct {
	Datadir string

	BTCMoveCheckEnable  bool
	BTCMoveCheckLogsURL string

	// 数据库重建模式
	DatabaseVersionRebuildMode bool
}

func NewEmptyChainStateConfig() *ChainStateConfig {
	cnf := &ChainStateConfig{}
	return cnf
}

func NewChainStateConfig(cnffile *sys.Inicnf) *ChainStateConfig {
	cnf := NewEmptyChainStateConfig()

	cnf.Datadir = path.Join(cnffile.MustDataDir(), "chainstate")

	// 验证比特币转移
	sec1 := cnffile.Section("btcmovecheck")
	cnf.BTCMoveCheckEnable = sec1.Key("enable").MustBool(false)
	cnf.BTCMoveCheckLogsURL = sec1.Key("logs_url").MustString("")

	//fmt.Println(cnf.BTCMoveCheckEnable, cnf.BTCMoveCheckLogsURL)

	return cnf
}
