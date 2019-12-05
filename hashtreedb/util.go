package hashtreedb

import "os"

///////////////////////////////////////////

func ReverseHashOrder(hash []byte) []byte {
	var length = len(hash)
	var hsdt = make([]byte, length)
	copy(hsdt, hash)
	for i := 0; i < length/2; i++ {
		hsdt[i], hsdt[length-i-1] = hsdt[length-i-1], hsdt[i]
	}
	return hsdt
}



//
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

