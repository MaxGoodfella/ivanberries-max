package logic

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"ivanberries-max/internal/cache"
	"ivanberries-max/internal/model"
	"ivanberries-max/internal/repository"
	"ivanberries-max/internal/service/validation/util"
	"ivanberries-max/internal/service/validation/validator"
	"log"
	"time"
)

type CategoryService struct {
	repo        *repository.CategoryRepository
	redisClient *cache.RedisClient
}

func NewCategoryService(repo *repository.CategoryRepository, redisClient *cache.RedisClient) *CategoryService {
	return &CategoryService{repo: repo, redisClient: redisClient}
}

func (s *CategoryService) GetByID(id uuid.UUID) (*model.Category, error) {
	cacheKey := fmt.Sprintf("category:%s", id.String())
	var category *model.Category

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
			return nil, util.ErrCategoryNotFound
		}
		return nil, err
	}

	categoryJSON, _ := json.Marshal(category)
	if err := s.redisClient.Set(cacheKey, string(categoryJSON), 10*time.Minute); err != nil {
		log.Printf("error setting category cache: %v", err)
	}

	return category, nil
}

func (s *CategoryService) GetAll() ([]model.Category, error) {
	cacheKey := "categories"
	var categories []model.Category

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

func (s *CategoryService) Create(category *model.Category) error {
	if err := validator.ValidateCategory(category); err != nil {
		return err
	}

	if err := s.repo.Create(category); err != nil {
		return util.ErrCategoryCreationFailed
	}

	return nil
}

func (s *CategoryService) Update(category *model.Category) (*model.Category, error) {
	if category.ID == uuid.Nil {
		return nil, util.ErrInvalidCategoryID
	}

	if err := validator.ValidateCategory(category); err != nil {
		return nil, err
	}

	_, err := s.repo.GetByID(category.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, util.ErrCategoryNotFound
		}
		return nil, err
	}

	if err := s.repo.Update(category); err != nil {
		return nil, util.ErrCategoryUpdateFailed
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
		return util.ErrInvalidCategoryID
	}

	if err := s.repo.Delete(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return util.ErrCategoryNotFound
		}
		return util.ErrCategoryDeleteFailed
	}

	cacheKey := fmt.Sprintf("category:%s", id.String())
	if err := s.redisClient.Delete(cacheKey); err != nil {
		log.Printf("Error deleting category cache: %v", err)
	}

	return nil
}
