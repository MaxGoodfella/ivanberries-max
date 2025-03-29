package services

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"ivanberries-max/internal/models"
	"ivanberries-max/internal/repositories"
	"math"
	"strconv"
	"strings"
)

type ProductService struct {
	repo *repositories.ProductRepository
}

func NewProductService(repo *repositories.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetProductByID(id uuid.UUID) (*models.Product, error) {
	return s.repo.GetByID(id)
}

func (s *ProductService) CreateProduct(product *models.Product) error {
	if strings.TrimSpace(product.Name) == "" {
		return errors.New("name cannot be null, empty, or just spaces")
	}

	if !nameRegexp.MatchString(product.Name) {
		return errors.New("invalid product name format")
	}

	if !containsLetter.MatchString(product.Name) {
		return errors.New("name must contain at least one letter")
	}

	if !descriptionRegexp.MatchString(product.Description) {
		return errors.New("invalid product description format")
	}

	if product.Price <= 0 {
		return errors.New("price must be greater than zero")
	}

	priceStr := strconv.FormatFloat(product.Price, 'f', -1, 64)

	if !priceRegexp.MatchString(priceStr) {
		return errors.New("price can contain a maximum of 2 decimal places")
	}

	return s.repo.Create(product)
}

func (s *ProductService) UpdateProduct(id uuid.UUID, updates map[string]interface{}) (*models.Product, error) {

	allowedFields := map[string]bool{
		"name":        true,
		"description": true,
		"price":       true,
		"category_id": true,
	}

	for field := range updates {
		if !allowedFields[field] {
			return nil, fmt.Errorf("invalid field: %s", field)
		}
	}

	_, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	if name, ok := updates["name"].(string); ok {
		if strings.TrimSpace(name) == "" {
			return nil, errors.New("name cannot be null, empty or just spaces")
		}
		if !nameRegexp.MatchString(name) {
			return nil, errors.New("invalid product name format")
		}
		if !containsLetter.MatchString(name) {
			return nil, errors.New("name must contain at least one letter")
		}
	}

	if description, ok := updates["description"].(string); ok {
		if !descriptionRegexp.MatchString(description) {
			return nil, errors.New("invalid product description format")
		}
		//if !containsLetter.MatchString(description) {
		//	return nil, errors.New("description must contain at least one letter")
		//}
	}

	if price, ok := updates["price"].(float64); ok {
		if price <= 0 {
			return nil, errors.New("price must be greater than 0")
		}

		priceStr := fmt.Sprintf("%.2f", price)
		if priceStr != fmt.Sprintf("%.2f", math.Floor(price)) {
			return nil, errors.New("price can contain a maximum of 2 decimal places")
		}

	}

	if categoryID, ok := updates["category_id"].(string); ok {
		var _ models.Category
		if err := s.repo.CheckCategoryExists(categoryID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, gorm.ErrRecordNotFound
			}
			return nil, err
		}
	}

	err = s.repo.Update(id, updates)
	if err != nil {
		return nil, err
	}

	updatedProduct, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return updatedProduct, nil
}

func (s *ProductService) Delete(id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("invalid product ID")
	}

	return s.repo.Delete(id)
}

func (s *ProductService) GetProducts() ([]models.Product, error) {
	return s.repo.GetAll()
}
