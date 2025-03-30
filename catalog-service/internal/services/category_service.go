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

type CategoryService struct {
	repo *repositories.CategoryRepository
}

func NewCategoryService(repo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) GetCategoryByID(id uuid.UUID) (*models.Category, error) {
	category, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utilities.ErrCategoryNotFound
		}
		return nil, err
	}
	return category, nil
}

func (s *CategoryService) GetCategories() ([]models.Category, error) {
	return s.repo.GetAll()
}

func (s *CategoryService) CreateCategory(category *models.Category) error {
	if err := validators.ValidateCategory(category); err != nil {
		return err
	}

	if err := s.repo.Create(category); err != nil {
		return utilities.ErrCategoryCreationFailed
	}

	return nil
}

func (s *CategoryService) UpdateCategory(category *models.Category) (*models.Category, error) {
	if category.ID == uuid.Nil {
		return nil, utilities.ErrInvalidCategoryID
	}

	if err := validators.ValidateCategory(category); err != nil {
		return nil, err
	}

	_, err := s.repo.GetByID(category.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utilities.ErrCategoryNotFound
		}
		return nil, err
	}

	if err := s.repo.Update(category); err != nil {
		return nil, utilities.ErrCategoryUpdateFailed
	}

	return s.repo.GetByID(category.ID)
}

func (s *CategoryService) Delete(id uuid.UUID) error {
	if id == uuid.Nil {
		return utilities.ErrInvalidCategoryID
	}

	if err := s.repo.Delete(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utilities.ErrCategoryNotFound
		}
		return utilities.ErrCategoryDeleteFailed
	}

	return nil
}
