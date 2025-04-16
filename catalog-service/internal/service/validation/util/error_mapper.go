package util

import (
	"errors"
	"gorm.io/gorm"
	"net/http"
)

func GetHTTPStatusCode(err error) int {
	switch {
	case
		errors.Is(err, ErrInvalidCategoryID),
		errors.Is(err, ErrCategoryNameEmpty),
		errors.Is(err, ErrCategoryNameInvalid),
		errors.Is(err, ErrCategoryNameNoLetter),

		errors.Is(err, ErrInvalidProductID),
		errors.Is(err, ErrProductNameEmpty),
		errors.Is(err, ErrProductNameInvalid),
		errors.Is(err, ErrProductNameNoLetter),
		errors.Is(err, ErrProductDescriptionInvalid),
		errors.Is(err, ErrProductPriceInvalid),
		errors.Is(err, ErrInvalidField):
		return http.StatusBadRequest

	case
		errors.Is(err, ErrCategoryNotFound),
		errors.Is(err, gorm.ErrRecordNotFound),

		errors.Is(err, ErrProductNotFound),
		errors.Is(err, gorm.ErrRecordNotFound):
		return http.StatusNotFound

	case
		errors.Is(err, ErrCategoryCreationFailed),
		errors.Is(err, ErrCategoryUpdateFailed),
		errors.Is(err, ErrCategoryDeleteFailed),

		errors.Is(err, ErrProductCreationFailed),
		errors.Is(err, ErrProductUpdateFailed),
		errors.Is(err, ErrProductDeleteFailed):
		return http.StatusInternalServerError

	default:
		return http.StatusInternalServerError
	}
}
