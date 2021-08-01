package blockstorev2

import (
	"github.com/hacash/core/sys"
	"path"
)

type BlockStoreConfig struct {
	Datadir string
}

func NewEmptyBlockStoreConfig() *BlockStoreConfig {
	cnf := &BlockStoreConfig{}
	return cnf
}

func NewBlockStoreConfig(cnffile *sys.Inicnf) *BlockStoreConfig {
	cnf := NewEmptyBlockStoreConfig()

	cnf.Datadir = path.Join(cnffile.MustDataDir(), "blockstore")

	return cnf
}
