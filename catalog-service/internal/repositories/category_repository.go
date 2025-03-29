package repositories

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"ivanberries-max/internal/models"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetByID(id uuid.UUID) (*models.Category, error) {
	var category models.Category

	if err := r.db.First(&category, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	return &category, nil
}

func (r *CategoryRepository) GetAll() ([]models.Category, error) {
	var categories []models.Category
	if err := r.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *CategoryRepository) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

func (r *CategoryRepository) CategoryExists(id uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&models.Category{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

func (r *CategoryRepository) Update(category *models.Category) error {
	result := r.db.Model(&models.Category{}).Where("id = ?", category.ID).Updates(category)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

func (r *CategoryRepository) Delete(id uuid.UUID) error {
	result := r.db.Where("id = ?", id).Delete(&models.Category{})
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}
