package util

import "errors"

var (
	ErrTokenMissing       = errors.New("authorization header missing")
	ErrTokenFormatInvalid = errors.New("invalid token format")
	ErrUserIDInvalid      = errors.New("invalid userID format")
	ErrUUIDInvalid        = errors.New("invalid UUID")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrBindJSONFailed     = errors.New("invalid request format")
)
