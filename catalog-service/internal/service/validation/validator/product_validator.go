package validator

import (
	"fmt"
	"ivanberries-max/internal/model"
	util2 "ivanberries-max/internal/service/validation/util"
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
		return util2.ErrProductNameEmpty
	}
	if !util2.NameRegexp.MatchString(product.Name) {
		return util2.ErrProductNameInvalid
	}
	if !util2.ContainsLetter.MatchString(product.Name) {
		return util2.ErrProductNameNoLetter
	}
	if !util2.DescriptionRegexp.MatchString(product.Description) {
		return util2.ErrProductDescriptionInvalid
	}
	if product.Price <= 0 {
		return util2.ErrProductPriceInvalid
	}
	priceStr := strconv.FormatFloat(product.Price, 'f', -1, 64)
	if !util2.PriceRegexp.MatchString(priceStr) {
		return util2.ErrProductPriceInvalid
	}
	return nil
}

func ValidateProductUpdates(updates map[string]interface{}) error {
	for field := range updates {
		if !allowedFields[field] {
			return fmt.Errorf("invalid field: %s", field)
		}
	}

	if name, ok := updates["name"].(string); ok {
		if strings.TrimSpace(name) == "" {
			return util2.ErrProductNameEmpty
		}
		if !util2.NameRegexp.MatchString(name) {
			return util2.ErrProductNameInvalid
		}
		if !util2.ContainsLetter.MatchString(name) {
			return util2.ErrProductNameNoLetter
		}
	}

	if description, ok := updates["description"].(string); ok {
		if !util2.DescriptionRegexp.MatchString(description) {
			return util2.ErrProductDescriptionInvalid
		}
	}

	if price, ok := updates["price"].(float64); ok {
		if price <= 0 {
			return util2.ErrProductPriceInvalid
		}
		priceStr := strconv.FormatFloat(price, 'f', -1, 64)
		if !util2.PriceRegexp.MatchString(priceStr) {
			return util2.ErrProductPriceInvalid
		}
	}

	return nil
}
