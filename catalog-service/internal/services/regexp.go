package services

import "regexp"

var (
	nameRegexp        = regexp.MustCompile(`^[a-zA-Z0-9](?:[a-zA-Z0-9 ]|[.,&#\-/])*[a-zA-Z0-9]$`)
	descriptionRegexp = regexp.MustCompile(`[a-zA-Z0-9]`)
	containsLetter    = regexp.MustCompile(`[a-zA-Z]`)
	priceRegexp       = regexp.MustCompile(`^\d+(\.\d{1,2})?$`)
)
