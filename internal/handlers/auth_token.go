package handlers

import (
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Token struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateJWT(username, secretKey string) string {

	var tokenClaim = Token{
		Username: username,
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

func VerifyJWT(bearerToken, secretKey string) (*jwt.Token, error) {

	token, err := jwt.ParseWithClaims(bearerToken, &Token{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	return token, err
}
