package hashing

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) ([]byte, error) {
	if bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost); err != nil {
		return nil, err
	} else {
		return bytes, nil
	}
}
func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func CheckPasswordHash(password string, hash string) (bool, error) {

	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) != nil {
		return false, errors.New("invalid username or password")
	} else {
		return true, nil
	}

}
