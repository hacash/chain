package biglogdb

import (
	"fmt"
	"testing"
)

func Test_t1(t *testing.T) {

	testdir := "/home/shiqiujie/Desktop/Hacash/go/src/github.com/hacash/chain/biglogdb/data1"

	dbcnf := NewBigLogDBConfig(testdir, 8, 300)
	dbcnf.LogHeadMaxSize = 5
	db, e := NewBigLogDB(dbcnf)
	if e != nil {
		panic(e)
	}

	//stoptr, err := db.Save([]byte("A2345678"), []byte(testdir))
	//if err != nil {
	//	fmt.Println(err)
	//}else{
	//	fmt.Println(stoptr.Filenum, stoptr.Fileseek, stoptr.Valsize)
	//}

	valuedts, e2 := db.Read([]byte("02345678"), 0)
	if e2 != nil {
		fmt.Println(e)
	}
	fmt.Println(string(valuedts))

}
