package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"users-service/internal/service/logic"
)

func RoleMiddleware(authService *logic.AuthService, requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token is missing"})
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
		if roleName == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "role not found"})
			c.Abort()
			return
		}

		// Проверяем, подходит ли роль
		isAllowed := false
		for _, r := range requiredRoles {
			if roleName == r {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID.String())
		c.Set("role", roleName)
		c.Next()
	}
}
