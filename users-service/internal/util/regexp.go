package util

import "regexp"

var (
	EmailRegexp    = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	PasswordRegexp = regexp.MustCompile(`^[a-zA-Z\d]{8,}$`)
	ContainsLetter = regexp.MustCompile(`[a-zA-Z]`)
)
