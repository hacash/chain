package statedomaindb

import "bytes"

// 查询实例

type QueryInstance struct {
	db *StateDomainDB

	inputkey  []byte
	domainkey []byte
}

func newQueryInstance(db *StateDomainDB, inkey []byte) (*QueryInstance, error) {
	// 给KEY 加上后缀
	keybuf := bytes.NewBuffer(inkey)
	keybuf.Write([]byte(db.config.KeyDomainName))
	// 创建
	ins, e := newQueryInstanceByRealUseKey(db, keybuf.Bytes())
	if e != nil {
		return nil, e
	}
	ins.inputkey = inkey
	return ins, nil
}

func newQueryInstanceByRealUseKey(db *StateDomainDB, realkey []byte) (*QueryInstance, error) {
	ins := &QueryInstance{
		db:        db,
		domainkey: realkey,
	}
	// 如果是内存数据库，则不打开本地文件
	if db.config.MemoryStorage {
		return ins, nil
	}
	// 如果是 level db, 则不打开文件
	if db.config.LevelDB {
		return ins, nil
	}

	panic("must use leveldb")

}

// 关闭
func (ins *QueryInstance) Destroy() {
	// 释放文件控制
	if !ins.db.config.MemoryStorage && !ins.db.config.LevelDB {
		// ins.db.releaseControlOfFile(ins)
	}
	// 清空数据
	ins.db = nil
	ins.inputkey = nil
	ins.domainkey = nil
}
