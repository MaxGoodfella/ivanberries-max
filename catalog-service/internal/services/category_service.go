package services

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"ivanberries-max/internal/models"
	"ivanberries-max/internal/repositories"
	"strings"
)

type CategoryService struct {
	repo *repositories.CategoryRepository
}

func NewCategoryService(repo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) GetCategoryByID(id uuid.UUID) (*models.Category, error) {
	return s.repo.GetByID(id)
}

func (s *CategoryService) CreateCategory(category *models.Category) error {
	name := strings.TrimSpace(category.Name)
	if name == "" {
		return errors.New("name cannot be null, empty, or just spaces")
	}

	if !nameRegexp.MatchString(name) {
		return errors.New("invalid category name format")
	}

	if !containsLetter.MatchString(name) {
		return errors.New("name must contain at least one letter")
	}

	return s.repo.Create(category)
}

func (s *CategoryService) UpdateCategory(category *models.Category) (*models.Category, error) {
	if category.ID == uuid.Nil {
		return nil, errors.New("invalid category ID")
	}

	name := strings.TrimSpace(category.Name)
	if name == "" {
		return nil, errors.New("name cannot be null, empty, or just spaces")
	}

	if !nameRegexp.MatchString(name) {
		return nil, errors.New("invalid category name format")
	}

	if !containsLetter.MatchString(name) {
		return nil, errors.New("name must contain at least one letter")
	}

	_, err := s.repo.GetByID(category.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	err = s.repo.Update(category)
	if err != nil {
		return nil, err
	}

	updatedCategory, err := s.repo.GetByID(category.ID)
	if err != nil {
		return nil, err
	}

	return updatedCategory, nil
}

func (s *CategoryService) Delete(id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("invalid category ID")
	}

	return s.repo.Delete(id)
}

func (s *CategoryService) GetCategories() ([]models.Category, error) {
	return s.repo.GetAll()
}
