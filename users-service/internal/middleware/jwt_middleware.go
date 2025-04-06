package middleware

import (
	"net/http"
	"strings"
	"users-service/internal/service/logic"

	"github.com/gin-gonic/gin"
)

func JWTMiddleware(authService *logic.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token missing"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})
			c.Abort()
			return
		}

		claims, err := authService.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		roleName := authService.GetRoleName(claims.RoleID)

		c.Set("userID", claims.UserID.String())
		c.Set("role", roleName)
		c.Next()
	}
}
