package randstring

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateRandomString генерирует случайную строку длиной n
func GenerateRandomString(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
