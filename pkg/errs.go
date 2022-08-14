package pkg

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
	ErrUserUnauthorized  = NewErr("user not authorized")

	ErrInvalidOrderNumber   = NewErr("invalid order number")
	ErrOrderAlreadyExists   = NewErr("order number already exists")
	ErrUserAlreadyOrderedIt = NewErr("user already ordered this order")

	ErrPaymentNotAvailable = NewErr("insufficient funds for payment")
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

	case
		ErrUserAlreadyExists,
		ErrOrderAlreadyExists:
		return http.StatusConflict

	case ErrEmptyAuthData:
		return http.StatusBadRequest

	case
		ErrUserNotFound,
		ErrUserUnauthorized:
		return http.StatusUnauthorized

	case ErrInvalidOrderNumber:
		return http.StatusUnprocessableEntity

	case ErrUserAlreadyOrderedIt:
		return http.StatusOK

	case ErrPaymentNotAvailable:
		return http.StatusPaymentRequired

	default:
		return http.StatusInternalServerError
	}
}
