package validator

import (
	"catalog-service/internal/model"
	"catalog-service/internal/service/validation/util"
	"strings"
)

func ValidateCategory(category *model.Category) error {
	name := strings.TrimSpace(category.Name)
	if name == "" {
		return util.ErrCategoryNameEmpty
	}

	if !util.NameRegexp.MatchString(name) {
		return util.ErrCategoryNameInvalid
	}

	if !util.ContainsLetter.MatchString(name) {
		return util.ErrCategoryNameNoLetter
	}

	return nil
}
