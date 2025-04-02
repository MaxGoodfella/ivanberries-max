package util

import "regexp"

var (
	NameRegexp        = regexp.MustCompile(`^[a-zA-Z0-9](?:[a-zA-Z0-9 ]|[.,&#\-/])*[a-zA-Z0-9]$`)
	DescriptionRegexp = regexp.MustCompile(`[a-zA-Z0-9]`)
	ContainsLetter    = regexp.MustCompile(`[a-zA-Z]`)
	PriceRegexp       = regexp.MustCompile(`^\d+(\.\d{1,2})?$`)
)
