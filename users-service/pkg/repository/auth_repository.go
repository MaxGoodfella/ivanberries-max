package repository

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
	"users-service/pkg/model"
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

	var existingToken model.RefreshToken
	err = r.db.Where("user_id = ?", parsedUserID).First(&existingToken).Error
	if err == nil {
		existingToken.Token = refreshToken
		existingToken.ExpiresAt = expiration
		return r.db.Save(&existingToken).Error
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
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

func (r *AuthRepository) GetPermissionsForRole(roleID int) ([]model.Permission, error) {
	var permissions []model.Permission

	// Получаем разрешения для роли
	err := r.db.Table("permissions").
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role_id = ?", roleID).
		Find(&permissions).Error

	if err != nil {
		return nil, err
	}

	return permissions, nil
}
