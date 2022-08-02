package auth

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

func VerifyJWT(bearerToken string, user market.User) (*jwt.Token, error) {

	token, err := jwt.ParseWithClaims(bearerToken, &TokenJWT{}, func(token *jwt.Token) (interface{}, error) {
		return user.Password, nil
	})

	return token, err
}

func GenerateJWT(user market.User) string {

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

func GeneratePasswordHash(password string) (string, error) {

	if len(password) < 1 {
		return "", market.ErrEmptyAuthData
	}

	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil)), nil
}
