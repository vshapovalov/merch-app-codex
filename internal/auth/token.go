package auth

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateToken produces a cryptographically secure random token string.
func GenerateToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}
