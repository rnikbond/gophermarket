package pkg

import (
	"crypto/sha256"
	"encoding/hex"
)

func GeneratePasswordHash(password, salt string) (string, error) {

	if len(password) < 1 {
		return ``, ErrEmptyAuthData
	}

	hash := sha256.New()
	hash.Write([]byte(password + salt))
	return hex.EncodeToString(hash.Sum(nil)), nil
}
