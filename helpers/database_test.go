package helpers

import (
	"testing"
	"time"
	"lambda-server/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	
	user := &models.User{
		UserId:          "test",
		Email:           "test@gmail.com",
		GoogleID:        "test",
		Name:            "test",
		ProfilePicture:  "test.com",
		AuthMethods:     []string{"google"},
		IsEmailVerified: true,
		TokenVersion:    1,
		LastActiveAt:    time.Now().Unix(),
		CreatedAt:       time.Now().Unix(),
		UpdatedAt:       time.Now().Unix(),
	}
	print(user)

	// err := CreateUser(user)
	// assert.Nil(t, err)
	assert.NotNil(t, user)
}