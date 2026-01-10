package helpers

import (
	"crypto/rand"
	"encoding/base64"
)

func IsEmptyString(s string) bool {
	if len(s) <= 0 {
		return true
	} else {
		return false
	}
}

func GenerateSessionToken(nBytes int) (string, error) {
	b := make([]byte, nBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
