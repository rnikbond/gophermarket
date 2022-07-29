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
	ErrEmptyLoginPassword = NewErr("login or password can no been empty")
	ErrUserAlreadyExists  = NewErr("login already registered")
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

	case ErrEmptyLoginPassword:
		return http.StatusBadRequest

	default:
		return http.StatusInternalServerError
	}
}
