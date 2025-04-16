package middleware

import (
	"github.com/MaxGoodfella/ivanberries-max/users-service/pkg/service/logic"
	"github.com/MaxGoodfella/ivanberries-max/users-service/pkg/util"
	"github.com/gin-gonic/gin"
	"strings"
)

func RoleMiddleware(authService *logic.AuthService, requiredRoles ...string) gin.HandlerFunc {
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
			c.JSON(util.GetHTTPStatusCode(util.ErrInvalidToken), gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		roleName := authService.GetRoleName(claims.RoleID)
		if roleName == "" {
			c.JSON(util.GetHTTPStatusCode(util.ErrRoleNotFound), gin.H{"error": util.ErrRoleNotFound.Error()})
			c.Abort()
			return
		}

		isAllowed := false
		for _, r := range requiredRoles {
			if roleName == r {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			c.JSON(util.GetHTTPStatusCode(util.ErrForbidden), gin.H{"error": util.ErrForbidden.Error()})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID.String())
		c.Set("role", roleName)
		c.Next()
	}
}
