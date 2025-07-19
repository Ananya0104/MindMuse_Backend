package routes

import (
	"lambda-server/handlers"
	"lambda-server/middlewares"

	"github.com/gin-gonic/gin"
)

// SetupJournalRoutes configures all journal-related routes
func SetupJournalRoutes(api *gin.RouterGroup) {
	journal := api.Group("/journals")
	{
		journal.GET("", middlewares.AuthMiddleware(), handlers.GetAllJournalEntries)
		journal.POST("", middlewares.AuthMiddleware(), handlers.CreateJournalEntry)
		journal.GET("/:journalId", middlewares.AuthMiddleware(), handlers.GetJournalEntry)
		journal.DELETE("/:journalId", middlewares.AuthMiddleware(), handlers.DeleteJournalEntry)
		journal.PUT("/:journalId", middlewares.AuthMiddleware(), handlers.UpdateJournalEntry)

		// Debug route to test query param extraction
		journal.GET("/test", func(c *gin.Context) {
			userId := c.Query("userId")
			c.JSON(200, gin.H{"userId": userId})
		})
	}
}
