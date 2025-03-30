package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"ivanberries-max/internal/cache"
	"ivanberries-max/internal/models"
	"ivanberries-max/internal/repositories"
	"ivanberries-max/internal/validation/utilities"
	"ivanberries-max/internal/validation/validators"
	"log"
	"time"
)

type CategoryService struct {
	repo        *repositories.CategoryRepository
	redisClient *cache.RedisClient
}

func NewCategoryService(repo *repositories.CategoryRepository, redisClient *cache.RedisClient) *CategoryService {
	return &CategoryService{repo: repo, redisClient: redisClient}
}

func (s *CategoryService) GetCategoryByID(id uuid.UUID) (*models.Category, error) {
	cacheKey := fmt.Sprintf("category:%s", id.String())
	var category *models.Category

	val, err := s.redisClient.Get(cacheKey)
	if err == nil && val != "" {
		err := json.Unmarshal([]byte(val), &category)
		if err == nil {
			return category, nil
		}
	}

	category, err = s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utilities.ErrCategoryNotFound
		}
		return nil, err
	}

	categoryJSON, _ := json.Marshal(category)
	if err := s.redisClient.Set(cacheKey, string(categoryJSON), 10*time.Minute); err != nil {
		log.Printf("error setting category cache: %v", err)
	}

	return category, nil
}

func (s *CategoryService) GetCategories() ([]models.Category, error) {
	cacheKey := "categories"
	var categories []models.Category

	val, err := s.redisClient.Get(cacheKey)
	if err == nil && val != "" {
		err := json.Unmarshal([]byte(val), &categories)
		if err == nil {
			return categories, nil
		}
	}

	categories, err = s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	categoriesJSON, _ := json.Marshal(categories)
	if err := s.redisClient.Set(cacheKey, string(categoriesJSON), 10*time.Minute); err != nil {
		log.Printf("Error setting categories cache: %v", err)
	}

	return categories, nil
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

	cacheKey := fmt.Sprintf("category:%s", category.ID.String())
	categoryJSON, _ := json.Marshal(category)
	if err := s.redisClient.Set(cacheKey, string(categoryJSON), 10*time.Minute); err != nil {
		log.Printf("Error updating category cache: %v", err)
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

	cacheKey := fmt.Sprintf("category:%s", id.String())
	if err := s.redisClient.Delete(cacheKey); err != nil {
		log.Printf("Error deleting category cache: %v", err)
	}

	return nil
}
