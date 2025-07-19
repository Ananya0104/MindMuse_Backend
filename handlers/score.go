package handlers

import (
	"context"
	"lambda-server/database"
	"lambda-server/models"
	"lambda-server/constants"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// SubmitMindMuseScore handles POST /score/submit
func SubmitMindMuseScore(c *gin.Context) {
	var req struct {
		Score float64 `json:"score" binding:"required"`
		Timestamp string `json:"timestamp" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	userId := ""
	if userIdVal, exists := c.Get(constants.ContextKeyUserId); exists {
		if id, ok := userIdVal.(string); ok {
			userId = id
		}
	}

	parsedTime, err := time.Parse(time.RFC3339, req.Timestamp)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid timestamp format"})
		return
	}

	scoreEntry := models.MindMuseScore{
		UserId:    userId,
		Score:     req.Score,
		Timestamp: parsedTime.Unix(),
	}

	ctx := context.Background()
	if err := database.CreateMindMuseScoreEntry(ctx, scoreEntry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit score", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Score submitted successfully"})
} 