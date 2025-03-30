package services

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"ivanberries-max/internal/models"
	"ivanberries-max/internal/repositories"
	"ivanberries-max/internal/validation/utilities"
	"ivanberries-max/internal/validation/validators"
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

func (s *ProductService) GetProducts() ([]models.Product, error) {
	return s.repo.GetAll()
}

func (s *ProductService) CreateProduct(product *models.Product) error {
	if err := validators.ValidateProduct(product); err != nil {
		return err
	}
	return s.repo.Create(product)
}

func (s *ProductService) UpdateProduct(id uuid.UUID, updates map[string]interface{}) (*models.Product, error) {
	if err := validators.ValidateProductUpdates(updates); err != nil {
		return nil, err
	}

	_, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utilities.ErrProductNotFound
		}
		return nil, err
	}

	if categoryID, ok := updates["category_id"].(string); ok {
		if err := s.repo.CheckCategoryExists(categoryID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, utilities.ErrCategoryNotFound
			}
			return nil, err
		}
	}

	err = s.repo.Update(id, updates)
	if err != nil {
		return nil, utilities.ErrProductUpdateFailed
	}

	return s.repo.GetByID(id)
}

func (s *ProductService) Delete(id uuid.UUID) error {
	if id == uuid.Nil {
		return utilities.ErrInvalidProductID
	}

	return s.repo.Delete(id)
}
