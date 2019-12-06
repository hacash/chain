package chainstore

type ChainStoreConfig struct {
	absdir string
}

func NewChainStoreConfig(absdir string) *ChainStoreConfig {
	cnf := &ChainStoreConfig{
		absdir: absdir,
	}
	return cnf
}
