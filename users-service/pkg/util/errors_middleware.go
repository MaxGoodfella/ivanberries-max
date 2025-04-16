package util

import "errors"

var (
	// JWT Middleware errors
	ErrBlacklistedToken = errors.New("blacklisted token")
	ErrInvalidToken     = errors.New("invalid token")

	// Role Middleware errors
	ErrRoleNotFound = errors.New("role not found")
	ErrForbidden    = errors.New("forbidden")
)
