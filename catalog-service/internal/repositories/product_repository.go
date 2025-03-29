package repositories

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"ivanberries-max/internal/models"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetByID(id uuid.UUID) (*models.Product, error) {
	var product models.Product

	if err := r.db.First(&product, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepository) GetAll() ([]models.Product, error) {
	var products []models.Product
	if err := r.db.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *ProductRepository) Update(id uuid.UUID, updates map[string]interface{}) error {
	result := r.db.Model(&models.Product{}).Where("id = ?", id).Updates(updates)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

func (r *ProductRepository) Delete(id uuid.UUID) error {
	result := r.db.Where("id = ?", id).Delete(&models.Product{})
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

func (r *ProductRepository) CheckCategoryExists(categoryID string) error {
	var category models.Category
	result := r.db.Where("id = ?", categoryID).First(&category)
	return result.Error
}
