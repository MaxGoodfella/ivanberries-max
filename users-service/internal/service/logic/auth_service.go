package logic

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
	"users-service/internal/cache"
	"users-service/internal/model"
	"users-service/internal/repository"
	"users-service/internal/service/validation"
	"users-service/internal/util"
)

type AuthService struct {
	repo        *repository.AuthRepository
	jwtSecret   string
	redisClient *cache.RedisClient
}

func NewAuthService(repo *repository.AuthRepository, jwtSecret string, redisClient *cache.RedisClient) *AuthService {
	return &AuthService{
		repo:        repo,
		jwtSecret:   jwtSecret,
		redisClient: redisClient,
	}
}

func (s *AuthService) Register(email, password, roleID string) (*model.User, error) {
	if err := validation.ValidateUser(email, password, roleID); err != nil {
		return nil, err
	}

	_, err := s.repo.GetByEmail(email)
	if err == nil {
		return nil, util.ErrUserEmailAlreadyRegistered
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	rID, err := strconv.ParseUint(roleID, 10, 0)
	if err != nil {
		return nil, util.ErrUserRoleIDInvalid
	}

	user := &model.User{
		ID:             uuid.New(),
		Email:          email,
		HashedPassword: string(hashedPassword),
		RoleID:         uint(rID),
		IsActive:       true,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Authenticate(email, password string) (string, string, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return "", "", util.ErrUserEmailOrPasswordInvalid
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
		return "", "", util.ErrUserEmailOrPasswordInvalid
	}

	accessToken, err := s.GenerateToken(user.ID.String(), strconv.Itoa(int(user.RoleID)), 15*time.Minute)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := s.GenerateToken(user.ID.String(), strconv.Itoa(int(user.RoleID)), 7*24*time.Hour)
	if err != nil {
		return "", "", err
	}

	err = s.redisClient.StoreRefreshToken(user.ID.String(), refreshToken, 7*24*time.Hour)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) RefreshTokens(refreshToken string) (string, string, error) {
	claims, err := s.ParseToken(refreshToken)
	if err != nil {
		return "", "", util.ErrTokenInvalid
	}

	storedToken, err := s.redisClient.GetRefreshToken(claims.UserID.String())
	if err != nil || storedToken != refreshToken {
		return "", "", util.ErrTokenInvalidOrExpired
	}

	accessToken, err := s.GenerateToken(claims.UserID.String(), strconv.Itoa(int(claims.RoleID)), 15*time.Minute)
	if err != nil {
		return "", "", err
	}
	newRefreshToken, err := s.GenerateToken(claims.UserID.String(), strconv.Itoa(int(claims.RoleID)), 7*24*time.Hour)
	if err != nil {
		return "", "", err
	}

	if err := s.redisClient.DeleteKey("refresh:" + claims.UserID.String()); err != nil {
		return "", "", err
	}
	if err := s.redisClient.StoreRefreshToken(claims.UserID.String(), newRefreshToken, 7*24*time.Hour); err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}

func (s *AuthService) GetByID(userID string) (*model.User, error) {
	return s.repo.GetByID(userID)
}

func (s *AuthService) GetRoleName(roleID uint) string {
	role, err := s.repo.FindRoleByID(roleID)
	if err != nil {
		return ""
	}
	return role.Name
}

func (s *AuthService) GenerateToken(userID, roleID string, duration time.Duration) (string, error) {
	uID, err := uuid.Parse(userID)
	if err != nil {
		return "", err
	}

	rID, err := strconv.ParseUint(roleID, 10, 0)
	if err != nil {
		fmt.Println("Error:", err)
		return "", nil
	}

	claims := &model.Claims{
		UserID: uID,
		RoleID: uint(rID),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *AuthService) ParseToken(tokenString string) (*model.Claims, error) {
	blacklisted, err := s.redisClient.IsTokenBlacklisted(tokenString)
	if err != nil {
		return nil, err
	}
	if blacklisted {
		return nil, util.ErrTokenBlacklisted
	}

	token, err := jwt.ParseWithClaims(tokenString, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, util.ErrSigningMethodInvalid
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*model.Claims)
	if !ok || !token.Valid {
		return nil, util.ErrTokenInvalid
	}

	return claims, nil
}

func (s *AuthService) Logout(token string) error {
	claims, err := s.ParseToken(token)
	if err != nil {
		return err
	}

	expiry := time.Until(claims.ExpiresAt.Time)
	return s.redisClient.BlacklistAccessToken(token, expiry)
}

func (s *AuthService) IsBlacklisted(token string) (bool, error) {
	return s.redisClient.IsTokenBlacklisted(token)
}
