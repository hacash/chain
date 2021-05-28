package hashtreedb

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"testing"
)

const (
	TestDir = "/media/yangjie/500GB/Hacash/src/github.com/hacash/chain/data1/"
)

func Test_gc(t *testing.T) {

	TestDir := "/home/shiqiujie/Desktop/Hacash/go/src/github.com/hacash/chain/hashtreedb/data1"

	cnf := NewHashTreeDBConfig(TestDir, 4, 4)
	//cnf.FileDividePartitionLevel = 2
	//cnf.KeepDeleteMark = true
	db := NewHashTreeDB(cnf)

	key1 := []byte{1, 1, 1, 1}
	key2 := []byte{17, 17, 17, 17}
	key3 := []byte{33, 33, 33, 33}

	fmt.Println(key1, key2, key3)

	kv := key2

	ins, _ := db.CreateNewQueryInstance(kv)
	ins.Save(kv)

	//ins.Delete()

	//fmt.Println( ins.Find() )
	//e3 := ins1.Delete()
	//if e3 != nil {
	//	fmt.Println(e3)
	//	return
	//}

	ins.Destroy()

}

func Test_store_list_t5(t *testing.T) {

	stokeys := [][]byte{
		//[]byte{0,0,0,0},
		//[]byte{0,0,0,1},
		//[]byte{0,1,0,2},
		[]byte{1, 1, 1, 3},
		//[]byte{1,0,1,3},
		//[]byte{0,1,1,3},
		//[]byte{1,1,1,3},
		[]byte{1, 1, 1, 1},
	}
	//stokey1 := []byte{1,7,8,9}
	cnf := NewHashTreeDBConfig(TestDir, 4, 4)
	//cnf.FileDividePartitionLevel = 2
	db := NewHashTreeDB(cnf)

	for _, key := range stokeys {

		fmt.Println("- - - - - - - -", key)

		ins1, _ := db.CreateNewQueryInstance(key)
		e1 := ins1.Save(key)
		if e1 != nil {
			fmt.Println(e1)
		}
		ins1.Destroy()

		ins2, _ := db.CreateNewQueryInstance(key)
		fddts2, e2 := ins2.Find()
		if e2 != nil {
			fmt.Println(e2)
		}
		fmt.Println(fddts2)
		ins2.Destroy()
	}

}

func Test_store_hash_t4(t *testing.T) {

	cnf := NewHashTreeDBConfig(TestDir, 32, 32)
	//cnf.FileDividePartitionLevel = 2
	db := NewHashTreeDB(cnf)

	curkey := make([]byte, 32)
	for i0 := 0; i0 < 1000000; i0++ {
		rand.Read(curkey)
		ins1, e0 := db.CreateNewQueryInstance(curkey)
		if e0 != nil {
			panic(e0)
		}
		e2 := ins1.Save(curkey)
		if e2 != nil {
			panic(e2)
		}
		ins1.Destroy()

		ins2, e3 := db.CreateNewQueryInstance(curkey)
		if e3 != nil {
			panic(e3)
		}
		fddts, e4 := ins2.Find()
		if e4 != nil {
			panic(e4)
		}
		ins2.Destroy()
		if bytes.Compare(fddts, curkey) != 0 {
			panic(fmt.Errorf("curkey[0:4] is error"))
		}
	}

	//
	//for i0:=uint8(0); i0<255; i0++ {
	//	for i1:=uint8(0); i1<10; i1++ {
	//		for i2:=uint8(0); i2<1; i2++ {
	//			for i3:=uint8(0); i3<1; i3++ {
	//				curkey := []byte{i0,i1,i2,i3}
	//				kkk := sha3.Sum256(curkey)
	//				ins1, _ := db.CreateNewQueryInstance(kkk[:])
	//				ins1.Save(curkey)
	//				ins1.Destroy()
	//			}
	//		}
	//	}
	//}

}

func Test_store_much_t3(t *testing.T) {

	cnf := NewHashTreeDBConfig(TestDir, 4, 4)
	//cnf.FileDividePartitionLevel = 2
	db := NewHashTreeDB(cnf)

	for i0 := uint8(0); i0 < 4; i0++ {
		for i1 := uint8(0); i1 < 4; i1++ {
			for i2 := uint8(0); i2 < 3; i2++ {
				for i3 := uint8(0); i3 < 3; i3++ {
					curkey := []byte{i0, i1, i2, i3}
					fmt.Println(curkey)
					ins1, e1 := db.CreateNewQueryInstance(curkey)
					if e1 != nil {
						panic(e1)
					}
					e2 := ins1.Save(curkey)
					if e2 != nil {
						panic(e2)
					}
					ins1.Destroy()

					ins2, e3 := db.CreateNewQueryInstance(curkey)
					if e3 != nil {
						panic(e3)
					}
					value, e4 := ins2.Find()
					if e4 != nil {
						panic(e4)
					}
					if bytes.Compare(value, curkey) != 0 {
						fmt.Println("...........", curkey, value)
						//panic("not value")
					} else {
						fmt.Println("===========", curkey, value)
					}
					ins2.Destroy()

				}
			}
		}
	}

}

func Test_store_new_t2(t *testing.T) {

	stokey1 := []byte{0, 0, 0, 0}
	stokey2 := []byte{0, 0, 0, 1}
	stokey3 := []byte{0, 0, 0, 2}
	stokey4 := []byte{0, 0, 0, 3}
	fmt.Println(stokey1, stokey2, stokey3, stokey4)
	//stokey1 := []byte{1,7,8,9}
	cnf := NewHashTreeDBConfig(TestDir, 4, 4)
	//cnf.FileDividePartitionLevel = 2
	db := NewHashTreeDB(cnf)

	ins1, _ := db.CreateNewQueryInstance(stokey1)
	ins1.Save([]byte("aaaa"))
	fddts1, _ := ins1.Find()
	fmt.Println(string(fddts1))
	ins1.Destroy()

	ins2, _ := db.CreateNewQueryInstance(stokey2)
	e2 := ins2.Save([]byte("VVVV"))
	fmt.Println(e2)
	fddts2, _ := ins2.Find()
	fmt.Println(string(fddts2))
	ins2.Destroy()

	//
	//ins3, _ := db.CreateNewQueryInstance(stokey3)
	//ins3.Save([]byte("BBBB"))
	//ins3.Destroy()

	//ins1.Save([]byte("aaaa"))
	//ins1.Save([]byte("BBBB"))
	//ins1.Save([]byte("bbbb"))
	//ins1.Save([]byte("VVVV"))
	//fddts, _ := ins1.Find()
	//fmt.Println(string(fddts))
	//ins1.Destroy()

}

func Test_create_query_ins_t1(t *testing.T) {

	key1, _ := hex.DecodeString("12a1633cafcc01ebfb6d78e39f687a1f0995c62fc95f51ead10a02ee0be551b5")
	fmt.Println(len(key1), key1)

	cnf := NewHashTreeDBConfig(TestDir, 80, 32)
	//cnf.FileDividePartitionLevel = 3
	db := NewHashTreeDB(cnf)

	//hash := db.convertKeyToHash(key1)

	//fpath, fk, usehash := db.locateTargetFilePath(hash)
	//
	//fmt.Println(fpath)
	//fmt.Println(fk)
	//fmt.Println(usehash)

	qins, err := db.CreateNewQueryInstance(key1)
	if err != nil {
		panic(err)
	}

	//fmt.Println(qins.filePath)

	// 关闭查询
	qins.Destroy()

	qins_1, _ := db.CreateNewQueryInstance(key1)
	qins_1.Destroy()

}
