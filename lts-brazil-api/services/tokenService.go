package services

import (
	"crypto/sha256"
	"encoding/hex"
)

func GenerateToken(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	token := hex.EncodeToString((hash.Sum((nil))))
	return token
}
