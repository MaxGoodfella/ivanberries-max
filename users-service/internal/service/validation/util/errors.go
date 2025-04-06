package util

import "errors"

var (
	ErrUserEmailEmpty    = errors.New("email cannot be null, empty, or just spaces")
	ErrUserEmailInvalid  = errors.New("invalid user email format")
	ErrUserEmailNoLetter = errors.New("email must contain at least one letter")

	ErrUserPasswordEmpty   = errors.New("password cannot be null, empty, or just spaces")
	ErrUserPasswordInvalid = errors.New("invalid user password format, password requires at least 8 characters, " +
		"including letters and digits")
	ErrUserEmailAlreadyRegistered = errors.New("email is already registered")

	ErrUserRoleIDInvalid = errors.New("invalid role ID")

	ErrUserEmailOrPasswordInvalid = errors.New("invalid email or password")

	ErrTokenInvalidOrExpired = errors.New("invalid or expired refresh token")
	ErrTokenInvalid          = errors.New("invalid refresh token")
	ErrSigningMethodInvalid  = errors.New("invalid signing method")
)
