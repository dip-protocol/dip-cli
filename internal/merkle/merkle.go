package merkle

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

// Hash computes a SHA256 hash of raw data and returns hex encoding.
func Hash(data []byte) string {

	digest := sha256.Sum256(data)

	return hex.EncodeToString(digest[:])
}

// Combine computes the Merkle parent of two child hashes.
//
// parent = SHA256(left_hash_bytes || right_hash_bytes)
func Combine(left string, right string) string {

	leftBytes, err := hex.DecodeString(left)
	if err != nil {
		panic(errors.New("invalid left hash"))
	}

	rightBytes, err := hex.DecodeString(right)
	if err != nil {
		panic(errors.New("invalid right hash"))
	}

	data := append(leftBytes, rightBytes...)

	digest := sha256.Sum256(data)

	return hex.EncodeToString(digest[:])
}