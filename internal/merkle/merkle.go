package merkle

import (
	"crypto/sha256"
	"encoding/hex"
)

func Hash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

func Combine(a string, b string) string {
	data := []byte(a + b)
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}