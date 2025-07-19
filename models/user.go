package models

// User represents a user item stored in DynamoDB
type User struct {
	UserId            string       `json:"userId" dynamodbav:"userId"`
	Name              string       `json:"name,omitempty" dynamodbav:"name,omitempty"`
	Username          string       `json:"username,omitempty" dynamodbav:"username,omitempty"`
	Email             string       `json:"email,omitempty" dynamodbav:"email,omitempty"`
	CountryCode       string       `json:"countryCode,omitempty" dynamodbav:"countryCode,omitempty"`
	Phone             string       `json:"phone,omitempty" dynamodbav:"phone,omitempty"`
	GoogleID          string       `json:"googleId,omitempty" dynamodbav:"googleId,omitempty"`
	PasswordHash      string       `json:"passwordHash,omitempty" dynamodbav:"passwordHash,omitempty"`
	AuthMethods       []string     `json:"authMethods" dynamodbav:"authMethods"` // ["email", "phone", "google"]
	IsEmailVerified   bool         `json:"isEmailVerified" dynamodbav:"isEmailVerified"`
	IsPhoneVerified   bool         `json:"isPhoneVerified" dynamodbav:"isPhoneVerified"`
	RefreshToken      string       `json:"refreshToken,omitempty" dynamodbav:"refreshToken,omitempty"`
	TokenVersion      int          `json:"tokenVersion" dynamodbav:"tokenVersion"`
	LastActiveAt      int64        `json:"lastActiveAt" dynamodbav:"lastActiveAt"`
	CreatedAt         int64        `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt         int64        `json:"updatedAt" dynamodbav:"updatedAt"`
	ProfilePicture    string       `json:"profilePicture,omitempty" dynamodbav:"profilePicture,omitempty"`
	Dob               string       `json:"dob,omitempty" dynamodbav:"dob,omitempty"` // date of birth
	EmergencyContacts [3]Emergency `json:"emergencyContacts,omitempty" dynamodbav:"emergencyContacts,omitempty"`
	// Password reset fields
	PasswordResetToken     string `json:"passwordResetToken,omitempty" dynamodbav:"passwordResetToken,omitempty"`
	PasswordResetExpiresAt int64  `json:"passwordResetExpiresAt,omitempty" dynamodbav:"passwordResetExpiresAt,omitempty"`
}

// LoginData represents login request data
type LoginData struct {
	Email    *string `json:"email" validate:"required,email"`
	Password *string `json:"password" validate:"required"`
}

// EmailChangePasswordRequestData represents request data for password reset
type EmailChangePasswordRequestData struct {
	Email *string `json:"email" validate:"required,email"`
}

// PasswordChangeRequestData represents password change data
type PasswordChangeRequestData struct {
	Password *string `json:"password" validate:"required"`
}

// UserUpdateRequest represents fields that can be updated in a user's profile
// All fields are optional; only provided fields will be updated
// Password should be plain text and will be hashed in the handler
// If updating email/phone, verification logic should be handled separately
type UserUpdateRequest struct {
	Name           *string `json:"name,omitempty"`
	Username       *string `json:"username,omitempty"`
	Email          *string `json:"email,omitempty"`
	CountryCode    *string `json:"countryCode,omitempty"`
	Phone          *string `json:"phone,omitempty"`
	Password       *string `json:"password,omitempty"`
	ProfilePicture *string `json:"profilePicture,omitempty"`
	Dob            *string `json:"dob,omitempty"`
}
