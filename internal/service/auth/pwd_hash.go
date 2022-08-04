package auth

import (
	"crypto/sha256"
	"encoding/hex"

	"gophermarket/pkg"
)

func (s Auth) GeneratePasswordHash(password string) (string, error) {

	if len(password) < 1 {
		return ``, pkg.ErrEmptyAuthData
	}

	hash := sha256.New()
	hash.Write([]byte(password + s.passwordSalt))
	return hex.EncodeToString(hash.Sum(nil)), nil
}
