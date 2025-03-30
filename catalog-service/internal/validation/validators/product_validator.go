package validators

import (
	"fmt"
	"ivanberries-max/internal/models"
	"ivanberries-max/internal/validation/utilities"
	"strconv"
	"strings"
)

var allowedFields = map[string]bool{
	"name":        true,
	"description": true,
	"price":       true,
	"category_id": true,
}

func ValidateProduct(product *models.Product) error {
	if strings.TrimSpace(product.Name) == "" {
		return utilities.ErrProductNameEmpty
	}
	if !utilities.NameRegexp.MatchString(product.Name) {
		return utilities.ErrProductNameInvalid
	}
	if !utilities.ContainsLetter.MatchString(product.Name) {
		return utilities.ErrProductNameNoLetter
	}
	if !utilities.DescriptionRegexp.MatchString(product.Description) {
		return utilities.ErrProductDescriptionInvalid
	}
	if product.Price <= 0 {
		return utilities.ErrProductPriceInvalid
	}
	priceStr := strconv.FormatFloat(product.Price, 'f', -1, 64)
	if !utilities.PriceRegexp.MatchString(priceStr) {
		return utilities.ErrProductPriceInvalid
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
			return utilities.ErrProductNameEmpty
		}
		if !utilities.NameRegexp.MatchString(name) {
			return utilities.ErrProductNameInvalid
		}
		if !utilities.ContainsLetter.MatchString(name) {
			return utilities.ErrProductNameNoLetter
		}
	}

	if description, ok := updates["description"].(string); ok {
		if !utilities.DescriptionRegexp.MatchString(description) {
			return utilities.ErrProductDescriptionInvalid
		}
	}

	if price, ok := updates["price"].(float64); ok {
		if price <= 0 {
			return utilities.ErrProductPriceInvalid
		}
		priceStr := strconv.FormatFloat(price, 'f', -1, 64)
		if !utilities.PriceRegexp.MatchString(priceStr) {
			return utilities.ErrProductPriceInvalid
		}
	}

	return nil
}
