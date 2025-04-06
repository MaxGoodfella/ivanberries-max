package repository

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
	"users-service/internal/model"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) GetByID(userID string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *AuthRepository) FindRoleByID(roleID uint) (*model.Role, error) {
	var role model.Role
	if err := r.db.First(&role, "id = ?", roleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &role, nil
}

func (r *AuthRepository) StoreRefreshToken(userID string, refreshToken string, expiration time.Time) error {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	refreshTokenRecord := model.RefreshToken{
		UserID:    parsedUserID,
		Token:     refreshToken,
		ExpiresAt: expiration,
	}

	return r.db.Create(&refreshTokenRecord).Error
}

func (r *AuthRepository) GetRefreshToken(token string) (*model.RefreshToken, error) {
	var refreshToken model.RefreshToken
	if err := r.db.Where("token = ?", token).First(&refreshToken).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &refreshToken, nil
}

func (r *AuthRepository) DeleteRefreshToken(token string) error {
	return r.db.Where("token = ?", token).Delete(&model.RefreshToken{}).Error
}
