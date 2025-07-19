package handlers

import (
	"errors"
	"fmt"

	"net/http"
	"time"

	"lambda-server/helpers"
	"lambda-server/models"
	"lambda-server/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func HandleLogin(c *gin.Context) {
	var authReq models.AuthRequest
	if err := c.ShouldBindJSON(&authReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
		})
		return
	}

	var user *models.User
	var err error

	switch authReq.AuthType {
	case "email":
		user, err = authenticateEmail(authReq.Credentials["email"], authReq.Credentials["password"])
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid authentication type",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// Update last active time
	user.LastActiveAt = time.Now().Unix()
	if err := helpers.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update user",
		})
		return
	}

	// Generate tokens
	tokens, err := helpers.GenerateTokenPair(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to generate tokens",
		})
		return
	}

	removeSensitiveInformationFromUser(user)

	helpers.SetAuthResponse(c, user, tokens)
}

// handleRegister processes registration requests
func HandleRegister(c *gin.Context) {
	var authReq models.AuthRequest
	if err := c.ShouldBindJSON(&authReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
		})
		return
	}

	var user *models.User
	var err error

	switch authReq.AuthType {
	case "email":
		dob := authReq.Credentials["dob"]
		err = utils.ValidateDOBFormat(dob)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		user, err = registerEmail(
			authReq.Credentials["email"],
			authReq.Credentials["password"],
			authReq.Credentials["name"],
			authReq.Credentials["phone"],
			authReq.Credentials["countryCode"],
			dob,
		)
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid authentication type",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// Generate tokens
	tokens, err := helpers.GenerateTokenPair(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to generate tokens",
		})
		return
	}

	removeSensitiveInformationFromUser(user)

	helpers.SetAuthResponse(c, user, tokens)
}


// handleRefresh processes token refresh requests
func HandleRefresh(c *gin.Context) {
	var refreshReq models.RefreshRequest
	if err := c.ShouldBindJSON(&refreshReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
		})
		return
	}

	tokens, err := helpers.RefreshTokens(refreshReq.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	helpers.SetRefreshTokensResponse(c, tokens)
}

// handleLogout processes logout requests
func HandleLogout(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "models.User context not found",
		})
		return
	}

	u := user.(*models.User)

	// Invalidate all tokens by incrementing token version
	u.TokenVersion++
	u.RefreshToken = ""
	if err := helpers.UpdateUser(u); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to logout",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logged out successfully",
	})
}

// handleGetProfile returns current user profile
func HandleGetProfile(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "models.User context not found",
		})
		return
	}

	u := user.(*models.User)
	removeSensitiveInformationFromUser(u)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user":    u,
	})
}

// UpdateCurrentUser handles PATCH /auth/me to update the current user's profile
func UpdateCurrentUser(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "User not found in context"})
		return
	}
	u := user.(*models.User)

	var req models.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request body", "details": err.Error()})
		return
	}

	updated := false

	if req.Name != nil {
		u.Name = *req.Name
		updated = true
	}
	if req.Username != nil {
		u.Username = *req.Username
		updated = true
	}
	if req.Email != nil {
		u.Email = *req.Email
		updated = true
	}
	if req.CountryCode != nil {
		u.CountryCode = *req.CountryCode
		updated = true
	}
	if req.Phone != nil {
		u.Phone = *req.Phone
		updated = true
	}
	if req.ProfilePicture != nil {
		u.ProfilePicture = *req.ProfilePicture
		updated = true
	}
	if req.Dob != nil {
		u.Dob = *req.Dob
		updated = true
	}

	if req.Password != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to hash password"})
			return
		}
		u.PasswordHash = string(hash)
		updated = true
	}

	if !updated {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "No valid fields to update"})
		return
	}

	if err := helpers.UpdateUser(u); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to update user"})
		return
	}

	removeSensitiveInformationFromUser(u)
	c.JSON(http.StatusOK, gin.H{"success": true, "user": u, "message": "Profile updated successfully"})
}

// DeleteCurrentUser handles DELETE /auth/me to delete the current user's account
func DeleteCurrentUser(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "User not found in context"})
		return
	}
	u := user.(*models.User)
	if err := helpers.DeleteUser(u.UserId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to delete user", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Account deleted successfully"})
}

// HandleForgotPassword handles the forgot password request
func HandleForgotPassword(c *gin.Context) {
	var req models.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request body"})
		return
	}
	user, err := helpers.GetUserByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "If the email exists, a reset link will be sent."})
		return
	}
	token := utils.GeneratePasswordResetToken()
	expiresAt := time.Now().Add(1 * time.Hour).Unix()
	user.PasswordResetToken = token
	user.PasswordResetExpiresAt = expiresAt
	helpers.UpdateUser(user)
	// Return token for frontend to send email
	resetLink := "https://godaiwellness.com/reset-password?token=" + token
	fmt.Println("Password reset link:", resetLink)
	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"message":   "Reset token generated successfully",
		"token":     token,
		"email":     req.Email,
		"resetLink": resetLink,
	})
}

// HandleResetPassword handles the reset password request
func HandleResetPassword(c *gin.Context) {
	var req models.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Token == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request body"})
		return
	}
	// Find user by reset token
	user, err := helpers.FindUserByResetToken(req.Token)
	if err != nil || user.PasswordResetExpiresAt < time.Now().Unix() {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid or expired token"})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to hash password"})
		return
	}
	user.PasswordHash = string(hash)
	user.PasswordResetToken = ""
	user.PasswordResetExpiresAt = 0
	helpers.UpdateUser(user)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Password has been reset successfully"})
}

// authenticateEmail validates email/password credentials
func authenticateEmail(email, password string) (*models.User, error) {
	user, err := helpers.GetUserByEmail(email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}

// registerEmail creates a new user with email/password
func registerEmail(email, password, name, phone, countryCode, dob string) (*models.User, error) {
	// Check if user already exists
	if _, err := helpers.GetUserByEmail(email); err == nil {
		return nil, errors.New("user already exists with this email")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	userID := utils.GenerateUserID()
	user := &models.User{
		UserId:          userID,
		Email:           email,
		Name:            name,
		Phone:           phone,
		CountryCode:     countryCode,
		Dob:             dob,
		PasswordHash:    string(hashedPassword),
		AuthMethods:     []string{"email"},
		IsEmailVerified: false, // In production, send verification email
		TokenVersion:    1,
		LastActiveAt:    time.Now().Unix(),
		CreatedAt:       time.Now().Unix(),
		UpdatedAt:       time.Now().Unix(),
	}

	if err := helpers.CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}


func removeSensitiveInformationFromUser(user *models.User) {
	user.PasswordHash = ""
	user.RefreshToken = ""
}
