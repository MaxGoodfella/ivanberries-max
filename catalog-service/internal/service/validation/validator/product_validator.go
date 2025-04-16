package validator

import (
	"catalog-service/internal/model"
	"catalog-service/internal/service/validation/util"
	"strconv"
	"strings"
)

var allowedFields = map[string]bool{
	"name":        true,
	"description": true,
	"price":       true,
	"category_id": true,
}

func ValidateProduct(product *model.Product) error {
	if strings.TrimSpace(product.Name) == "" {
		return util.ErrProductNameEmpty
	}
	if !util.NameRegexp.MatchString(product.Name) {
		return util.ErrProductNameInvalid
	}
	if !util.ContainsLetter.MatchString(product.Name) {
		return util.ErrProductNameNoLetter
	}
	if !util.DescriptionRegexp.MatchString(product.Description) {
		return util.ErrProductDescriptionInvalid
	}
	if product.Price <= 0 {
		return util.ErrProductPriceInvalid
	}
	priceStr := strconv.FormatFloat(product.Price, 'f', -1, 64)
	if !util.PriceRegexp.MatchString(priceStr) {
		return util.ErrProductPriceInvalid
	}
	return nil
}

func ValidateProductUpdates(updates map[string]interface{}) error {
	for field := range updates {
		if !allowedFields[field] {
			return util.ErrInvalidField
		}
	}

	if name, ok := updates["name"].(string); ok {
		if strings.TrimSpace(name) == "" {
			return util.ErrProductNameEmpty
		}
		if !util.NameRegexp.MatchString(name) {
			return util.ErrProductNameInvalid
		}
		if !util.ContainsLetter.MatchString(name) {
			return util.ErrProductNameNoLetter
		}
	}

	if description, ok := updates["description"].(string); ok {
		if !util.DescriptionRegexp.MatchString(description) {
			return util.ErrProductDescriptionInvalid
		}
	}

	if price, ok := updates["price"].(float64); ok {
		if price <= 0 {
			return util.ErrProductPriceInvalid
		}
		priceStr := strconv.FormatFloat(price, 'f', -1, 64)
		if !util.PriceRegexp.MatchString(priceStr) {
			return util.ErrProductPriceInvalid
		}
	}

	return nil
}
