package handlers

import (
	"context"
	"lambda-server/constants"
	"lambda-server/database"
	"lambda-server/models"
	"lambda-server/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateJournalEntry handles POST /journals/
func CreateJournalEntry(c *gin.Context) {
	var req models.JournalCreateRequest
	var userId string
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}
	userId = c.Query(constants.QueryParamUserId)
	if userId == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Missing userId in query params",
		})
		return
	}

	ctx := context.Background()

	currentTime := time.Now()
	entry := models.Journal{
		UserId:    userId,
		CreatedAt: currentTime.Unix(),
		JournalID: utils.GenerateJournalID(),
		Title:     req.Title,
		Content:   req.Content,
		Date:      currentTime.Format("20060102"),
		UpdatedAt: currentTime.Unix(),
	}

	err := database.CreateJournalEntry(ctx, entry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to create journal entry",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.JournalResponse{
		Journal: entry,
		Message: "Journal entry created successfully",
	})
}

// GetJournalEntry handles GET /journals/:journalId?userId=...
func GetJournalEntry(c *gin.Context) {
	journalId := c.Param(constants.QueryParamJournalId)
	if journalId == "" {
		journalId = c.Query(constants.QueryParamJournalId)
	}
	if journalId == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Missing journalId in path or query params",
		})
		return
	}
	userId := c.Query(constants.QueryParamUserId)
	if userId == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Missing userId in query params",
		})
		return
	}
	ctx := context.Background()

	foundEntry, err := database.GetJournalByID(ctx, userId, journalId)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.JournalResponse{Journal: *foundEntry})
}

// DeleteJournalEntry handles DELETE /journals/:journalId?userId=...
func DeleteJournalEntry(c *gin.Context) {
	journalId := c.Param(constants.QueryParamJournalId)
	if journalId == "" {
		journalId = c.Query(constants.QueryParamJournalId)
	}
	if journalId == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Missing journalId in path or query params",
		})
		return
	}
	userId := c.Query(constants.QueryParamUserId)
	if userId == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Missing userId in query params",
		})
		return
	}
	ctx := context.Background()

	err := database.DeleteJournalEntry(ctx, userId, journalId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to delete journal entry",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Journal entry deleted successfully"})
}

// GetAllJournalEntries handles GET /journals?userId=...
func GetAllJournalEntries(c *gin.Context) {
	userId := c.Query(constants.QueryParamUserId)
	if userId == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Missing userId in query params",
		})
		return
	}
	ctx := context.Background()

	// For now, scan all and filter in handler (for demo/testing)
	entries, err := database.GetUserJournals(ctx, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to fetch journal entries",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.JournalListResponse{
		Journals: entries,
		Count:    len(entries),
	})
}

// UpdateJournalEntry handles PUT /journals/:journalId?userId=...
func UpdateJournalEntry(c *gin.Context) {
	journalId := c.Param(constants.QueryParamJournalId)
	if journalId == "" {
		journalId = c.Query(constants.QueryParamJournalId)
	}
	if journalId == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Missing journalId in path or query params",
		})
		return
	}
	userId := c.Query(constants.QueryParamUserId)
	if userId == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Missing userId in query params",
		})
		return
	}

	var updateData models.JournalUpdateRequest

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}
	ctx := context.Background()

	err := database.UpdateJournalEntry(ctx, userId, journalId, updateData.Title, updateData.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to update journal entry",
			Details: err.Error(),
		})
		return
	}

	// Fetch the updated entry to get the correct updatedAt value
	updatedEntry, err := database.GetJournalByID(ctx, userId, journalId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to fetch updated journal entry",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.JournalResponse{
		Journal: *updatedEntry,
		Message: "Journal entry updated successfully",
	})
}
