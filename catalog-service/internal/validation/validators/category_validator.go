package validators

import (
	"ivanberries-max/internal/models"
	"ivanberries-max/internal/validation/utilities"
	"strings"
)

func ValidateCategory(category *models.Category) error {
	name := strings.TrimSpace(category.Name)
	if name == "" {
		return utilities.ErrCategoryNameEmpty
	}

	if !utilities.NameRegexp.MatchString(name) {
		return utilities.ErrCategoryNameInvalid
	}

	if !utilities.ContainsLetter.MatchString(name) {
		return utilities.ErrCategoryNameNoLetter
	}

	return nil
}
