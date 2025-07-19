package helpers

import (
	"lambda-server/models"
	"lambda-server/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func SetAuthResponse(c *gin.Context, user *models.User, tokens *models.TokenPair) {
	c.JSON(http.StatusOK, models.AuthResponse{
		Success: true,
		User:    user,
		Tokens:  tokens,
	})
}

func SetRefreshTokensResponse(c *gin.Context, tokens *models.TokenPair) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"tokens":  tokens,
	})
}

// CreateNewUserFromGoogleUser creates a new user against a google user. And returns the user.
//
// The function performs the following operations:
// 1. Creates new user with Google auth
//
// Input Parameters:
//   - googleUser (*models.GoogleUser): The Google User Object containing google details of user
//
// Response:
//   - Success: Returns user object
//   - Error: Returns error message
func CreateNewUserFromGoogleUser(googleUser *models.GoogleUser) (*models.User, error) {
	userID := utils.GenerateUserID()
	user := &models.User{
		UserId:          userID,
		Email:           googleUser.Email,
		GoogleID:        googleUser.ID,
		Name:            googleUser.Name,
		ProfilePicture:  googleUser.Picture,
		AuthMethods:     []string{"google"},
		IsEmailVerified: googleUser.VerifiedEmail,
		TokenVersion:    1,
		LastActiveAt:    time.Now().Unix(),
		CreatedAt:       time.Now().Unix(),
		UpdatedAt:       time.Now().Unix(),
	}

	if err := CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func AddGoogleSigninToUser(user *models.User, googleUser *models.GoogleUser) {
	user.GoogleID = googleUser.ID
	user.ProfilePicture = googleUser.Picture
	user.IsEmailVerified = googleUser.VerifiedEmail
	// Update name if provided by Google and different
	if googleUser.Name != "" && googleUser.Name != user.Name {
		user.Name = googleUser.Name
	}
	// Add 'google' to AuthMethods if not already present
	found := false
	for _, method := range user.AuthMethods {
		if method == "google" {
			found = true
			break
		}
	}
	if !found {
		user.AuthMethods = append(user.AuthMethods, "google")
	}
}
