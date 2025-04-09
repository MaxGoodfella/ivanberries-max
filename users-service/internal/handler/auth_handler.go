package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"users-service/internal/model"
	"users-service/internal/service/logic"
	"users-service/internal/util"
)

type AuthHandler struct {
	authService *logic.AuthService
}

func NewAuthHandler(authService *logic.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(util.GetHTTPStatusCode(util.ErrBindJSONFailed), gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.Register(req.Email, req.Password, strconv.Itoa(int(req.RoleID)))
	if err != nil {
		c.JSON(util.GetHTTPStatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         user.ID,
		"email":      user.Email,
		"role_id":    user.RoleID,
		"is_active":  user.IsActive,
		"created_at": user.CreatedAt,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(util.GetHTTPStatusCode(util.ErrBindJSONFailed), gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, err := h.authService.Authenticate(req.Email, req.Password)
	if err != nil {
		c.JSON(util.GetHTTPStatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *AuthHandler) Me(c *gin.Context) {
	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(util.GetHTTPStatusCode(util.ErrUnauthorized), gin.H{"error": util.ErrUnauthorized.Error()})
		return
	}

	userIDStr, ok := userIDRaw.(string)
	if !ok {
		c.JSON(util.GetHTTPStatusCode(util.ErrUserIDInvalid), gin.H{"error": util.ErrUserIDInvalid.Error()})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(util.GetHTTPStatusCode(util.ErrUUIDInvalid), gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.GetByID(userID.String())
	if err != nil {
		c.JSON(util.GetHTTPStatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         user.ID,
		"email":      user.Email,
		"role_id":    user.RoleID,
		"is_active":  user.IsActive,
		"created_at": user.CreatedAt,
	})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req model.TokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(util.GetHTTPStatusCode(util.ErrBindJSONFailed), gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, err := h.authService.RefreshTokens(req.RefreshToken)
	if err != nil {
		c.JSON(util.GetHTTPStatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(util.GetHTTPStatusCode(util.ErrTokenMissing), gin.H{"error": util.ErrTokenMissing.Error()})
		return
	}

	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	err := h.authService.Logout(tokenString)
	if err != nil {
		c.JSON(util.GetHTTPStatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully logged out"})
}
