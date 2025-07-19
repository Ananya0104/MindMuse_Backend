package models

import (
	jwt "github.com/dgrijalva/jwt-go"
)

type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiresIn"`
	TokenType    string `json:"tokenType"`
}

type JWTClaims struct {
	UserID       string `json:"userId"`
	TokenVersion int    `json:"tokenVersion"`
	TokenType    string `json:"tokenType"` // "access" or "refresh"
	jwt.StandardClaims
}

// GoogleUser represents the Google OAuth user information
type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

// GoogleTokenInfo represents Google token validation response
type GoogleTokenInfo struct {
	Audience      string `json:"aud"`
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Locale        string `json:"locale"`
	ExpiresIn     string `json:"expires_in"`
}