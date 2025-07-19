package handlers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifyGoogleTokenFunc(t *testing.T) {
	googleToken := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImJiNDM0Njk1OTQ0NTE4MjAxNDhiMzM5YzU4OGFlZGUzMDUxMDM5MTkiLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20iLCJhenAiOiI2Mjg0MTY5NDYxMTctbW5rbWk5MTI2cnNpMHJvcGpvdGpkamJqZ2VnbGlnb2wuYXBwcy5nb29nbGV1c2VyY29udGVudC5jb20iLCJhdWQiOiI2Mjg0MTY5NDYxMTctbW5rbWk5MTI2cnNpMHJvcGpvdGpkamJqZ2VnbGlnb2wuYXBwcy5nb29nbGV1c2VyY29udGVudC5jb20iLCJzdWIiOiIxMDM0MjYxNzAwOTY2ODc1MDM0MTYiLCJoZCI6InVuaXRlZHRlY2hsYWIuY29tIiwiZW1haWwiOiJrdWxkZWVwQHVuaXRlZHRlY2hsYWIuY29tIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsIm5iZiI6MTc0OTEyODg3MiwibmFtZSI6Ikt1bGRlZXAgSmFuZ2lyIiwicGljdHVyZSI6Imh0dHBzOi8vbGgzLmdvb2dsZXVzZXJjb250ZW50LmNvbS9hL0FDZzhvY0xad0h1akMxY1FtVUJJN2JUMTkzRGo3U05wQzVlWFA3RFV2VGY1VjVTMmsyN1dWMEk9czk2LWMiLCJnaXZlbl9uYW1lIjoiS3VsZGVlcCIsImZhbWlseV9uYW1lIjoiSmFuZ2lyIiwiaWF0IjoxNzQ5MTI5MTcyLCJleHAiOjE3NDkxMzI3NzIsImp0aSI6IjBjYjkyN2M2Mzg4MjM4YmIxZDlmMWY3N2NkNWY1OTYzYjQxZmY3YmIifQ.S4RBFvZ0JDTcM2nL7ouVZu2Nx_ciLGpUycTqvt0sahtIGpp2nuWrWJ24b8sv9fegWyXZFuYyvJ3ktvDeRsegAATGWWcvWh3tiK99ihxajajnr02TO-RO7EHHvszZ5DN9l2jPifN78n-mInnR5P0ko68Fhg4UO_wCvjBozZtteQSfq3KqH6tB6oCzjCPTogQh8ZUvHRiwmEsx2LV_WuqUJj7-5gkqPm29rI1qE5AMaXvXGFCztlvqEE3ds2GWD6dFMlDPw-aNfHLiUG6qMM5cIX9Hawf1PebRcmvalTm-vYRSWcC6PBBRsxWY1xdRaX3-p59BZavPm2G9uoJZVRRyBA"

	// user , err := VerifyGoogleTokenFunc(googleToken)
	// print(user)
	// assert.Nil(t, err)
	assert.NotEmpty(t, googleToken)
}

func TestRegisterProcess(t *testing.T) {

	// if _, err := helpers.GetUserByGoogleID("103426170096687503416"); err == nil {
	// 	print("user already exists with this Google account")
	// 	return
	// } else {
	// 	print("\nuser not exists with this Google account\n")
	// 	print(err.Error())
	// }

	// if _, err := helpers.GetUserByEmail("kuldeep@unitedtechlab.com"); err == nil {
	// 	print("user already exists with this email")
	// 	return
	// } else {
	// 	print("\nuser not exists with this email\n")
	// 	print(err.Error())
	// }

}

// func TestEmailSignupThenGoogleLoginAddsAuthMethod(t *testing.T) {
// 	gin.SetMode(gin.TestMode)
// 	r := gin.New()
// 	routes.SetupAuthRoutes(r.Group(""))

// 	// Use a unique email for each test run
// 	email := fmt.Sprintf("testuser_%d@example.com", time.Now().UnixNano())
// 	password := "testpassword"
// 	name := "Test User"
// 	dob := "01/11/2002"

// 	// Step 1: Register with email/password
// 	registerPayload := models.AuthRequest{
// 		AuthType: "email",
// 		Credentials: map[string]string{
// 			"email":       email,
// 			"password":    password,
// 			"name":        name,
// 			"dob":         dob,
// 			"phone":       "",
// 			"countryCode": "",
// 		},
// 	}
// 	registerBody, _ := json.Marshal(registerPayload)
// 	w := httptest.NewRecorder()
// 	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(registerBody))
// 	req.Header.Set("Content-Type", "application/json")
// 	r.ServeHTTP(w, req)
// 	assert.Equal(t, http.StatusOK, w.Code)

// 	var registerResp models.AuthResponse
// 	_ = json.Unmarshal(w.Body.Bytes(), &registerResp)
// 	assert.True(t, registerResp.Success)
// 	assert.Equal(t, email, registerResp.User.Email)
// 	assert.Contains(t, registerResp.User.AuthMethods, "email")

// 	// Step 2: Simulate Google login for the same email
// 	originalVerifyGoogleIDToken := handlers.VerifyGoogleIDToken
// 	defer func() { handlers.VerifyGoogleIDToken = originalVerifyGoogleIDToken }()
// 	handlers.VerifyGoogleIDToken = func(idTokenString, clientID string) (*models.GoogleUser, error) {
// 		return &models.GoogleUser{
// 			ID:            "google-id-123",
// 			Email:         email, // Use the same email as registration
// 			VerifiedEmail: true,
// 			Name:          "Google Name",
// 			Picture:       "http://example.com/pic.jpg",
// 		}, nil
// 	}

// 	googlePayload := models.GoogleAuthRequest{
// 		GoogleToken: "fake-google-token",
// 		Action:      "login",
// 	}
// 	googleBody, _ := json.Marshal(googlePayload)
// 	w2 := httptest.NewRecorder()
// 	req2, _ := http.NewRequest("POST", "/auth/google", bytes.NewBuffer(googleBody))
// 	req2.Header.Set("Content-Type", "application/json")
// 	r.ServeHTTP(w2, req2)
// 	assert.Equal(t, http.StatusOK, w2.Code)

// 	var googleResp models.AuthResponse
// 	_ = json.Unmarshal(w2.Body.Bytes(), &googleResp)
// 	assert.True(t, googleResp.Success)
// 	assert.Equal(t, email, googleResp.User.Email)
// 	assert.Contains(t, googleResp.User.AuthMethods, "email")
// 	assert.Contains(t, googleResp.User.AuthMethods, "google")
// }
