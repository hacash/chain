package chainstate

import (
	"github.com/hacash/core/sys"
	"path"
)

type ChainStateConfig struct {
	Datadir string

	SatoshiEnable         bool
	SatoshiBTCMoveLogsURL string
}

func NewEmptyChainStateConfig() *ChainStateConfig {
	cnf := &ChainStateConfig{}
	return cnf
}

func NewChainStateConfig(cnffile *sys.Inicnf) *ChainStateConfig {
	cnf := NewEmptyChainStateConfig()

	cnf.Datadir = path.Join(cnffile.MustDataDir(), "chainstate")

	// 验证比特币转移
	sec1 := cnffile.Section("satoshi")
	cnf.SatoshiEnable = sec1.Key("enable").MustBool(false)
	cnf.SatoshiBTCMoveLogsURL = sec1.Key("btcmovelogs_url").MustString("")

	//fmt.Println(cnf.SatoshiEnable, cnf.SatoshiBTCMoveLogsURL)

	return cnf
}
