package tinykvdb

import "crypto/md5"

func convertKeyToLen16Hash(key []byte) []byte {
	hash := md5.Sum(key)
	return hash[:]
}
