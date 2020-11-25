package hashtreedb

import "sync"

type lockFilePkgItem struct {
	count                  int
	lock                   *sync.Mutex
	targetFilePackageCache *TargetFilePackage
}
