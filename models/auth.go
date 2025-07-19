package models

// AuthRequest represents authentication request payload
type AuthRequest struct {
	AuthType    string            `json:"authType"` // "email", "phone", "google"
	Credentials map[string]string `json:"credentials"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	Success bool       `json:"success"`
	Message string     `json:"message,omitempty"`
	User    *User      `json:"user,omitempty"`
	Tokens  *TokenPair `json:"tokens,omitempty"`
}

// RefreshRequest represents token refresh request
type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

// GoogleAuthRequest represents Google OAuth request
type GoogleAuthRequest struct {
	GoogleToken string `json:"googleToken"`
	Action      string `json:"action"` // "login" or "register"
}

// ForgotPasswordRequest represents forgot password request
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ResetPasswordRequest represents reset password request
type ResetPasswordRequest struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
}
