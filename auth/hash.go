package auth

import (
	"crypto/sha256"
	"encoding/hex"
)

func hash(data []byte) []byte {
	if len(data) == 0 {
		return []byte{}
	}
	b := sha256.Sum256(data)
	return b[:]
}

func HashString(data string) string {
	return hex.EncodeToString(hash([]byte(data)))
}
