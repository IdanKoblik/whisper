package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashToken(rawToken string) string {
	hash := sha256.Sum256([]byte(rawToken))
	return hex.EncodeToString(hash[:])
}
