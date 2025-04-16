package util

import (
	"errors"
	"gorm.io/gorm"
	"net/http"
)

func GetHTTPStatusCode(err error) int {
	switch {
	// Service errors
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
		errors.Is(err, ErrTokenInvalid),
		errors.Is(err, ErrTokenBlacklisted):
		return http.StatusUnauthorized

	case errors.Is(err, gorm.ErrRecordNotFound):
		return http.StatusNotFound

	// Handler errors
	case errors.Is(err, ErrTokenMissing),
		errors.Is(err, ErrTokenFormatInvalid),
		errors.Is(err, ErrUUIDInvalid),
		errors.Is(err, ErrUserIDInvalid),
		errors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized

	case errors.Is(err, ErrBindJSONFailed):
		return http.StatusBadRequest

	// Middleware errors
	case errors.Is(err, ErrBlacklistedToken):
		return http.StatusUnauthorized
	case errors.Is(err, ErrInvalidToken):
		return http.StatusUnauthorized
	case errors.Is(err, ErrRoleNotFound):
		return http.StatusForbidden
	case errors.Is(err, ErrForbidden):
		return http.StatusForbidden

	default:
		return http.StatusInternalServerError
	}
}
