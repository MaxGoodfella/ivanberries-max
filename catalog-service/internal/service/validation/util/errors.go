package util

import "errors"

var (
	ErrInvalidCategoryID      = errors.New("invalid category ID")
	ErrCategoryNameEmpty      = errors.New("name cannot be null, empty, or just spaces")
	ErrCategoryNameInvalid    = errors.New("invalid category name format")
	ErrCategoryNameNoLetter   = errors.New("name must contain at least one letter")
	ErrCategoryNotFound       = errors.New("category not found")
	ErrCategoryCreationFailed = errors.New("failed to create category")
	ErrCategoryUpdateFailed   = errors.New("failed to update category")
	ErrCategoryDeleteFailed   = errors.New("failed to delete category")

	ErrInvalidProductID          = errors.New("invalid product ID")
	ErrProductNameEmpty          = errors.New("name cannot be null, empty, or just spaces")
	ErrProductNameInvalid        = errors.New("invalid product name format")
	ErrProductNameNoLetter       = errors.New("name must contain at least one letter")
	ErrProductDescriptionInvalid = errors.New("invalid product description format")
	ErrProductPriceInvalid       = errors.New("price must be greater than 0 and can contain a maximum of 2 decimal places")
	ErrProductNotFound           = errors.New("product not found")
	ErrProductCreationFailed     = errors.New("failed to create product")
	ErrProductUpdateFailed       = errors.New("failed to update product")
	ErrProductDeleteFailed       = errors.New("failed to delete product")
	ErrInvalidField              = errors.New("invalid field name")
)
