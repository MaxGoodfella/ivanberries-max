package middleware

import (
	"github.com/gin-gonic/gin"
	"strings"
	"users-service/pkg/service/logic"
	"users-service/pkg/util"
)

func PermissionMiddleware(authService *logic.AuthService, requiredPermissions ...string) gin.HandlerFunc {
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

		claims, err := authService.ParseToken(parts[1])
		if err != nil {
			c.JSON(util.GetHTTPStatusCode(util.ErrInvalidToken), gin.H{"error": util.ErrInvalidToken.Error()})
			c.Abort()
			return
		}

		roleName := authService.GetRoleName(claims.RoleID)
		if roleName == "" {
			c.JSON(util.GetHTTPStatusCode(util.ErrRoleNotFound), gin.H{"error": util.ErrRoleNotFound.Error()})
			c.Abort()
			return
		}

		hasPermission := false
		for _, permission := range requiredPermissions {
			if authService.HasPermission(int(claims.RoleID), permission) {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			c.JSON(util.GetHTTPStatusCode(util.ErrForbidden), gin.H{"error": util.ErrForbidden.Error()})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID.String())
		c.Set("role", roleName)
		c.Next()
	}
}
