package chainstate

type ChainStateConfig struct {
	Datadir string
}

func NewChainStateConfig(datadir string) *ChainStateConfig {
	cnf := &ChainStateConfig{
		Datadir: datadir,
	}
	return cnf
}
