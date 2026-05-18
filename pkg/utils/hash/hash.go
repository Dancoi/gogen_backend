package hash

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashToken хеширует токен SHA256
func HashToken(token string) string {
	h := sha256.New()
	h.Write([]byte(token))
	return hex.EncodeToString(h.Sum(nil))
}

// VerifyTokenHash проверяет токен с хешем
func VerifyTokenHash(token, tokenHash string) bool {
	return HashToken(token) == tokenHash
}
