package model

import (
	"github.com/golang-jwt/jwt/v5"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Email          string    `gorm:"unique;not null" json:"email"`
	HashedPassword string    `gorm:"not null" json:"-"`
	RoleID         uint      `gorm:"not null" json:"role_id"`
	IsActive       bool      `gorm:"default:true" json:"is_active"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type Role struct {
	ID   uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"unique;not null" json:"name"`
}

type Permission struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Code        string `gorm:"unique;not null" json:"code"`
	Description string `gorm:"not null" json:"description"`
}

type RolePermission struct {
	RoleID       uint `gorm:"not null;primaryKey" json:"role_id"`
	PermissionID uint `gorm:"not null;primaryKey" json:"permission_id"`
}

type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Token     string    `gorm:"unique;not null" json:"token"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
}

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	RoleID uint      `json:"role_id"`
	jwt.RegisteredClaims
}

type RegisterUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	RoleID   uint   `json:"role_id"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}
