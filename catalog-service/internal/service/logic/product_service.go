package logic

import (
	"catalog-service/internal/kafka"
	"catalog-service/internal/model"
	"catalog-service/internal/repository"
	"catalog-service/internal/service/validation/util"
	"catalog-service/internal/service/validation/validator"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
	"os"
)

type ProductService struct {
	repo     *repository.ProductRepository
	producer *kafka.Producer
}

func NewProductService(repo *repository.ProductRepository, producer *kafka.Producer) *ProductService {
	return &ProductService{repo: repo, producer: producer}
}

func (s *ProductService) GetByID(id uuid.UUID) (*model.Product, error) {
	return s.repo.GetByID(id)
}

func (s *ProductService) GetAll() ([]model.Product, error) {
	return s.repo.GetAll()
}

func (s *ProductService) Create(product *model.Product) error {
	if err := validator.ValidateProduct(product); err != nil {
		return err
	}
	return s.repo.Create(product)
}

func (s *ProductService) Update(id uuid.UUID, updates map[string]interface{}) (*model.Product, error) {
	if err := validator.ValidateProductUpdates(updates); err != nil {
		return nil, err
	}

	_, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, util.ErrProductNotFound
		}
		return nil, err
	}

	if categoryID, ok := updates["category_id"].(string); ok {
		if err := s.repo.CheckCategoryExists(categoryID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, util.ErrCategoryNotFound
			}
			return nil, err
		}
	}

	err = s.repo.Update(id, updates)
	if err != nil {
		log.Printf("Update failed: %v", err)
		return nil, util.ErrProductUpdateFailed
	}

	event := fmt.Sprintf(`{"event":"%s","product_id":"%s"}`, os.Getenv("KAFKA_EVENT_PRODUCT_UPDATED"), id)
	err = s.producer.SendMessage(id.String(), event)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
	} else {
		log.Println("Message successfully sent to Kafka")
	}

	return s.repo.GetByID(id)
}

func (s *ProductService) Delete(id uuid.UUID) error {
	if id == uuid.Nil {
		return util.ErrInvalidProductID
	}

	return s.repo.Delete(id)
}
