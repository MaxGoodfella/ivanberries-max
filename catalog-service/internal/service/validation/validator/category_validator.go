package validator

import (
	"ivanberries-max/internal/model"
	util2 "ivanberries-max/internal/service/validation/util"
	"strings"
)

func ValidateCategory(category *model.Category) error {
	name := strings.TrimSpace(category.Name)
	if name == "" {
		return util2.ErrCategoryNameEmpty
	}

	if !util2.NameRegexp.MatchString(name) {
		return util2.ErrCategoryNameInvalid
	}

	if !util2.ContainsLetter.MatchString(name) {
		return util2.ErrCategoryNameNoLetter
	}

	return nil
}
