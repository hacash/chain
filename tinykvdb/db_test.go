package tinykvdb

import (
	"fmt"
	"testing"
)

func Test_t1(t *testing.T) {

	testdir := "/home/shiqiujie/Desktop/Hacash/go/src/github.com/hacash/chain/tinykvdb/data1"

	kv, _ := NewTinyKVDB(testdir, false)
	kv.Set([]byte("abc"), []byte(testdir))
	kv.Set([]byte("123"), []byte("testdir"))

	//kv.Del([]byte("123"))

	val, _ := kv.Get([]byte("123"))
	fmt.Println(string(val))

}
