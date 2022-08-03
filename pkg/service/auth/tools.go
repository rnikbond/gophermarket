package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
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

func ParseCookie(r *http.Request) (string, error) {

	cookie, err := r.Cookie("token")
	if err != nil {
		fmt.Println("cookie not exists")
		return ``, market.ErrEmptyAuthData
	}

	bearerToken := cookie.Value
	token, err := VerifyJWT(bearerToken)

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return ``, market.ErrUserUnauthorized
		}

		return ``, market.ErrEmptyAuthData
	}

	if !token.Valid {
		return ``, market.ErrUserUnauthorized
	}

	user := token.Claims.(*TokenJWT)
	return user.Username, nil
}

func VerifyJWT(bearerToken string) (*jwt.Token, error) {

	token, err := jwt.ParseWithClaims(bearerToken, &TokenJWT{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
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
