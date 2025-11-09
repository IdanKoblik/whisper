package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestHashToken(t *testing.T) {
	t.Run("valid token produces correct SHA256 hash", func(t *testing.T) {
		input := "my-secret-token"
		expectedBytes := sha256.Sum256([]byte(input))
		expectedHash := hex.EncodeToString(expectedBytes[:])

		result := HashToken(input)
		assert.Equal(t, expectedHash, result)
	})

	t.Run("empty string produces correct SHA256 hash", func(t *testing.T) {
		input := ""
		expectedBytes := sha256.Sum256([]byte(input))
		expectedHash := hex.EncodeToString(expectedBytes[:])

		result := HashToken(input)
		assert.Equal(t, expectedHash, result)
	})

	t.Run("different inputs produce different hashes", func(t *testing.T) {
		hash1 := HashToken("token1")
		hash2 := HashToken("token2")
		assert.NotEqual(t, hash1, hash2)
	})
}
