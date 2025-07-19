package helpers

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"lambda-server/constants"
	"lambda-server/models"

	"github.com/dgrijalva/jwt-go"
)

var (
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
)

func GenerateTokenPair(user *models.User) (*models.TokenPair, error) {
	accessTokenExpiry := time.Now().Add(30 * time.Minute)
	refreshTokenExpiry := time.Now().Add(30 * 24 * time.Hour)

	// Generate access token
	accessClaims := &models.JWTClaims{
		UserID:       user.UserId,
		TokenVersion: user.TokenVersion,
		TokenType:    constants.TokenTypeAccess,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: accessTokenExpiry.Unix(),
			IssuedAt:  time.Now().Unix(),
			Subject:   user.UserId,
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(jwtSecret)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshClaims := &models.JWTClaims{
		UserID:       user.UserId,
		TokenVersion: user.TokenVersion,
		TokenType:    constants.TokenTypeRefresh,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: refreshTokenExpiry.Unix(),
			IssuedAt:  time.Now().Unix(),
			Subject:   user.UserId,
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(jwtSecret)
	if err != nil {
		return nil, err
	}

	// Update user's refresh token in database
	user.RefreshToken = refreshTokenString
	user.UpdatedAt = time.Now().Unix()
	if err := UpdateUser(user); err != nil {
		return nil, err
	}

	return &models.TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresIn:    accessTokenExpiry.Unix(),
		TokenType:    "Bearer",
	}, nil
}

// refreshTokens generates new token pair using refresh token
func RefreshTokens(refreshToken string) (*models.TokenPair, error) {
	// Validate refresh token
	claims, err := ValidateToken(refreshToken, constants.TokenTypeRefresh)
	if err != nil {
		return nil, err
	}

	// Get user and check token version
	user, err := GetUserByID(claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.TokenVersion != claims.TokenVersion || user.RefreshToken != refreshToken {
		return nil, errors.New("invalid refresh token")
	}

	// Check inactivity
	if err := CheckInactivity(user); err != nil {
		return nil, err
	}

	// Generate new token pair
	return GenerateTokenPair(user)
}

// validateToken validates JWT token
func ValidateToken(tokenString, tokenType string) (*models.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*models.JWTClaims); ok && token.Valid {
		if claims.TokenType != tokenType {
			return nil, errors.New("invalid token type")
		}
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// checkInactivity checks if user has been inactive for too long
func CheckInactivity(user *models.User) error {
	inactivityThreshold, _ := strconv.ParseInt(os.Getenv("INACTIVITY_THRESHOLD"), 10, 64)
	if inactivityThreshold == 0 {
		inactivityThreshold = 30 * 24 * 60 * 60 // 30 days default
	}

	if time.Now().Unix()-user.LastActiveAt > inactivityThreshold {
		// Invalidate tokens by incrementing version
		user.TokenVersion++
		user.RefreshToken = ""
		UpdateUser(user)
		return errors.New("session expired due to inactivity")
	}

	return nil
}
