package auth

import (
	"crypto/sha256"
	"encoding/hex"
)

func Hash(data []byte) []byte {
	if len(data) == 0 {
		return []byte{}
	}
	b := sha256.Sum256(data)
	return b[:]
}

func HashString(data string) string {
	return hex.EncodeToString(Hash([]byte(data)))
}
