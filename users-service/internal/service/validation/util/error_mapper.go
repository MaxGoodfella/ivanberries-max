package util

import (
	"errors"
	"gorm.io/gorm"
	"net/http"
)

func GetHTTPStatusCode(err error) int {
	switch {
	case errors.Is(err, ErrUserEmailEmpty),
		errors.Is(err, ErrUserEmailInvalid),
		errors.Is(err, ErrUserPasswordEmpty),
		errors.Is(err, ErrUserPasswordInvalid),
		errors.Is(err, ErrUserRoleIDInvalid),
		errors.Is(err, ErrSigningMethodInvalid):
		return http.StatusBadRequest

	case errors.Is(err, ErrUserEmailAlreadyRegistered):
		return http.StatusConflict

	case errors.Is(err, ErrUserEmailOrPasswordInvalid),
		errors.Is(err, ErrTokenInvalidOrExpired),
		errors.Is(err, ErrTokenInvalid):
		return http.StatusUnauthorized

	case errors.Is(err, gorm.ErrRecordNotFound):
		return http.StatusNotFound

	default:
		return http.StatusInternalServerError
	}
}
