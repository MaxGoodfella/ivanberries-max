package repository

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"ivanberries-max/internal/model"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetByID(id uuid.UUID) (*model.Category, error) {
	var category model.Category

	if err := r.db.First(&category, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	return &category, nil
}

func (r *CategoryRepository) GetAll() ([]model.Category, error) {
	var categories []model.Category
	if err := r.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *CategoryRepository) Create(category *model.Category) error {
	return r.db.Create(category).Error
}

func (r *CategoryRepository) CategoryExists(id uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&model.Category{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

func (r *CategoryRepository) Update(category *model.Category) error {
	result := r.db.Model(&model.Category{}).Where("id = ?", category.ID).Updates(category)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

func (r *CategoryRepository) Delete(id uuid.UUID) error {
	result := r.db.Where("id = ?", id).Delete(&model.Category{})
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}
