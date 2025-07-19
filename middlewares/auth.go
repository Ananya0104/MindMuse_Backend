package middlewares

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"lambda-server/constants"
	"lambda-server/helpers"
	"lambda-server/models"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT tokens
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// var tokensNeedRefresh = false
		claims, err := useAccessTokenToGetClaims(c)
		if err != nil {
			log.Println("Unable to use Access Token")
			// claims, err = useRefreshTokenToGetClaims(c)
			// if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Authorization header required or invalid",
			})
			c.Abort()
			return
			// }
			// tokensNeedRefresh = true
		}

		// Get user and validate token version
		user, err := helpers.GetUserByID(claims.UserID)
		if err != nil || user.TokenVersion != claims.TokenVersion {
			log.Println("Token was expired in version number or invalid", err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Token expired or invalid",
			})
			c.Abort()
			return
		}

		// Check inactivity
		if err := helpers.CheckInactivity(user); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": err.Error(),
			})
			c.Abort()
			return
		}

		// if tokensNeedRefresh {
		// 	tokens, err := helpers.GenerateTokenPair(user)
		// 	if err != nil {
		// 		c.JSON(http.StatusInternalServerError, gin.H{
		// 			"success": false,
		// 			"message": "Failed to generate new tokens",
		// 		})
		// 		c.Abort()
		// 		return
		// 	}
		// 	// Return new tokens in response headers
		// 	c.Header("X-New-Access-Token", tokens.AccessToken)
		// 	c.Header("X-New-Refresh-Token", tokens.RefreshToken)
		// }

		// Update last active time
		user.LastActiveAt = time.Now().Unix()
		helpers.UpdateUser(user)

		// Set user in context
		c.Set("user", user)
		c.Set("userId", user.UserId)
		c.Next()
	}
}

func useAccessTokenToGetClaims(c *gin.Context) (*models.JWTClaims, error) {
	tokenString, err := GetTokenFromAuthorizationHeader(c)
	if err != nil {
		println("AuthMiddleware: Failed to get accessToken from header: ", err.Error())
		return nil, errors.New("Access token absent")
	}
	claims, err := helpers.ValidateToken(tokenString, constants.TokenTypeAccess)
	if err != nil {
		println("Invalid Access Token: ", err.Error())
		return nil, errors.New("Access token invalid")
	}
	return claims, nil
}

func useRefreshTokenToGetClaims(c *gin.Context) (*models.JWTClaims, error) {
	tokenString, err := GetTokenFromAuthorizationHeader(c)
	if err != nil {
		println("AuthMiddleware: Failed to get refreshToken from header: ", err.Error())
		return nil, errors.New("Refresh token absent")
	}
	claims, err := helpers.ValidateToken(tokenString, constants.TokenTypeRefresh)
	if err != nil {
		println("Invalid Refresh Token: ", err.Error())
		return nil, errors.New("Refresh token invalid")
	}
	return claims, nil
}

func GetTokenFromAuthorizationHeader(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header missing")
	}

	parts := strings.Fields(authHeader)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", errors.New("invalid authorization header format")
	}

	return parts[1], nil
}
