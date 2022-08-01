package service

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"time"

	market "gophermarket/pkg"

	"github.com/dgrijalva/jwt-go"
)

const (
	secretKey = "secretKeyJWT"
)

type TokenJWT struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func verifyJWT(bearerToken string, user market.User) (*jwt.Token, error) {

	token, err := jwt.ParseWithClaims(bearerToken, &TokenJWT{}, func(token *jwt.Token) (interface{}, error) {
		return user.Password, nil
	})

	return token, err
}

func (s AuthService) generateJWT(user market.User) string {

	var tokenClaim = TokenJWT{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(10 * time.Minute).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaim)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		log.Fatal(err)
	}

	return tokenString
}

func (s AuthService) generatePasswordHash(user market.User) string {
	hash := sha256.New()
	hash.Write([]byte(user.Password + pwdSalt))
	return hex.EncodeToString(hash.Sum(nil))
}
