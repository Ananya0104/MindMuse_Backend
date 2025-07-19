package middlewares

import (
	"lambda-server/helpers"

	"github.com/gin-gonic/gin"
)

// SoftAuthMiddleware tries to authenticate but never aborts the request.
func SoftAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := useAccessTokenToGetClaims(c)
		if err == nil {
			user, err := helpers.GetUserByID(claims.UserID)
			if err == nil && user.TokenVersion == claims.TokenVersion {
				c.Set("user", user)
				c.Set("userId", user.UserId)
			}
		}
		// Always continue, even if not authenticated
		c.Next()
	}
}
