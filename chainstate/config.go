package chainstate

import (
	"github.com/hacash/core/sys"
	"path"
)

type ChainStateConfig struct {
	Datadir string
}

func NewEmptyChainStateConfig() *ChainStateConfig {
	cnf := &ChainStateConfig{}
	return cnf
}

func NewChainStateConfig(cnffile *sys.Inicnf) *ChainStateConfig {
	cnf := NewEmptyChainStateConfig()

	cnf.Datadir = path.Join(cnffile.MustDataDir(), "chainstore")

	return cnf
}
