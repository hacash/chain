package blockstorev3

import (
	"github.com/hacash/core/sys"
	"path"
)

type BlockStoreConfig struct {
	Datadir string

	// btc move
	BTCMoveCheckEnable    bool
	DownloadBTCMoveLogUrl string
}

func NewEmptyBlockStoreConfig() *BlockStoreConfig {
	cnf := &BlockStoreConfig{}
	return cnf
}

func NewBlockStoreConfig(cnffile *sys.Inicnf) *BlockStoreConfig {
	cnf := NewEmptyBlockStoreConfig()

	cnf.Datadir = path.Join(cnffile.MustDataDirWithVersion(), "blockstore")

	sec2 := cnffile.Section("btcmovecheck")
	cnf.BTCMoveCheckEnable = sec2.Key("enable").MustBool(false)
	cnf.DownloadBTCMoveLogUrl = sec2.Key("logs_url").MustString("")

	return cnf
}
