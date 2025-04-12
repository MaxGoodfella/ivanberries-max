package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"users-service/pkg/service/logic"
	"users-service/pkg/util"
)

func JWTMiddleware(authService *logic.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(util.GetHTTPStatusCode(util.ErrTokenMissing), gin.H{"error": util.ErrTokenMissing.Error()})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(util.GetHTTPStatusCode(util.ErrTokenFormatInvalid), gin.H{"error": util.ErrTokenFormatInvalid.Error()})
			c.Abort()
			return
		}

		token := parts[1]

		isBlacklisted, err := authService.IsBlacklisted(token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		if isBlacklisted {
			c.JSON(util.GetHTTPStatusCode(util.ErrBlacklistedToken), gin.H{"error": util.ErrBlacklistedToken.Error()})
			c.Abort()
			return
		}

		claims, err := authService.ParseToken(token)
		if err != nil {
			c.JSON(util.GetHTTPStatusCode(util.ErrInvalidToken), gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		roleName := authService.GetRoleName(claims.RoleID)

		c.Set("userID", claims.UserID.String())
		c.Set("role", roleName)
		c.Next()
	}
}
