package blockstorev3

import "bytes"

func keyfix(k []byte, suffix string) []byte {
	buf := bytes.NewBuffer(k)
	buf.Write([]byte(suffix))
	return buf.Bytes()
}
