package chainstatev3

import (
	"github.com/hacash/core/sys"
	"path"
)

type ChainStateConfig struct {
	Datadir string

	BTCMoveCheckEnable  bool
	BTCMoveCheckLogsURL string

	// Database rebuild mode
	DatabaseVersionRebuildMode bool
}

func NewEmptyChainStateConfig() *ChainStateConfig {
	cnf := &ChainStateConfig{}
	return cnf
}

func NewChainStateConfig(cnffile *sys.Inicnf) *ChainStateConfig {
	cnf := NewEmptyChainStateConfig()

	cnf.Datadir = path.Join(cnffile.MustDataDirWithVersion(), "chainstate")

	// Verify bitcoin transfer
	sec1 := cnffile.Section("btcmovecheck")
	cnf.BTCMoveCheckEnable = sec1.Key("enable").MustBool(false)
	cnf.BTCMoveCheckLogsURL = sec1.Key("logs_url").MustString("")

	//fmt.Println(cnf.BTCMoveCheckEnable, cnf.BTCMoveCheckLogsURL)

	return cnf
}
