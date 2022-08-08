package handlers

import (
	"log"
	"time"

	market "gophermarket/internal"

	"github.com/dgrijalva/jwt-go"
)

type Token struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func (h *Handler) GenerateJWT(user market.User) string {

	var tokenClaim = Token{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(10 * time.Minute).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaim)
	tokenString, err := token.SignedString([]byte(h.tokenKey))
	if err != nil {
		log.Fatal(err)
	}

	return tokenString
}

func (h *Handler) VerifyJWT(bearerToken string) (*jwt.Token, error) {

	token, err := jwt.ParseWithClaims(bearerToken, &Token{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.tokenKey), nil
	})

	return token, err
}
