package gophermarket

import (
	"errors"
	"net/http"
)

type ErrGM struct {
	Value string
}

func NewErr(s string) ErrGM {
	return ErrGM{Value: s}
}

func (es ErrGM) Error() string {
	return es.Value
}

var (
	ErrEmptyAuthData     = NewErr("empty login or password")
	ErrUserAlreadyExists = NewErr("login already registered")
	ErrUserNotFound      = NewErr("user not found")
)

var (
	ErrGenerateToken = NewErr("internal error generate token")
)

// ErrorHTTP - Преобразование ошибки Storage в HTTP код
func ErrorHTTP(err error) int {

	var serviceErr ErrGM
	if !errors.As(err, &serviceErr) {
		return http.StatusInternalServerError
	}

	switch serviceErr {

	case ErrUserAlreadyExists:
		return http.StatusConflict

	case ErrEmptyAuthData:
		return http.StatusBadRequest

	case ErrUserNotFound:
		return http.StatusUnauthorized

	default:
		return http.StatusInternalServerError
	}
}
