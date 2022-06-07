package hashtreedb

import "sync"

type MemoryStorageItem struct {
	IsDelete bool   // Tag deleted
	Value    []byte // data
}

type MemoryStorageDB struct {
	wlok  sync.Mutex
	Datas map[string]*MemoryStorageItem
}

func NewMemoryStorageDB() *MemoryStorageDB {
	return &MemoryStorageDB{
		Datas: make(map[string]*MemoryStorageItem),
	}
}

// clean
func (ms *MemoryStorageDB) Clean() {
	ms.wlok.Lock()
	defer ms.wlok.Unlock()
	ms.Datas = make(map[string]*MemoryStorageItem)
}

// len
func (ms *MemoryStorageDB) Len() int {
	return len(ms.Datas)
}

// save
func (ms *MemoryStorageDB) Save(realkey []byte, value []byte) {
	ms.wlok.Lock()
	defer ms.wlok.Unlock()
	if v, has := ms.Datas[string(realkey)]; has {
		v.IsDelete = false
		v.Value = value
	} else {
		ms.Datas[string(realkey)] = &MemoryStorageItem{
			IsDelete: false,
			Value:    value,
		}
	}
}

// exist
func (ms *MemoryStorageDB) Exist(realkey []byte) bool {
	if v, has := ms.Datas[string(realkey)]; has {
		if v.IsDelete {
			return false // deleted
		}
		return true // success
	} else {
		return false // not find
	}
}

// read
func (ms *MemoryStorageDB) Read(realkey []byte) ([]byte, bool) {
	if v, has := ms.Datas[string(realkey)]; has {
		if v.IsDelete {
			return nil, false // deleted
		}
		return v.Value, true
	} else {
		return nil, false // not find
	}
}

// delete
func (ms *MemoryStorageDB) Delete(realkey []byte) {
	ms.wlok.Lock()
	defer ms.wlok.Unlock()
	if v, has := ms.Datas[string(realkey)]; has {
		v.IsDelete = true // Delete tag
		v.Value = nil
	} else {
		ms.Datas[string(realkey)] = &MemoryStorageItem{
			IsDelete: true, // 删除标记
		}
	}
}
