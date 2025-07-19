package handlers

import (
	"lambda-server/constants"
	"lambda-server/database"
	"lambda-server/models"
	"lambda-server/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateOrUpdateEmergencyContactsForUser allows admin to set contacts for any user
func CreateEmergencyContacts(c *gin.Context) {
	userId := c.Param(constants.DynamoDbKeyUserId)
	if userId == constants.EMPTY_STRING {
		userId = c.Query(constants.QueryParamUserId)
	}
	if userId == constants.EMPTY_STRING {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is required"})
		return
	}

	var req struct {
		Contacts [3]models.Emergency `json:"contacts" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	// Convert array to slice for validation
	err := utils.ValidateEmergencyContacts(req.Contacts[:])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = database.SetEmergencyContacts(userId, req.Contacts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set emergency contacts", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"userId":   userId,
		"contacts": req.Contacts,
		"message":  "Emergency contacts set successfully",
	})
}

// GetEmergencyContacts fetches emergency contacts for a user
func GetEmergencyContacts(c *gin.Context) {
	userId := c.Param(constants.DynamoDbKeyUserId)
	if userId == constants.EMPTY_STRING {
		userId = c.Query(constants.QueryParamUserId)
	}
	if userId == constants.EMPTY_STRING {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is required"})
		return
	}

	contacts, err := database.GetEmergencyContacts(userId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Emergency contacts not found", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"userId":   userId,
		"contacts": contacts,
	})
}
