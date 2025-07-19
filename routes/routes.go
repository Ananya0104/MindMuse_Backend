package routes

import (
	"lambda-server/handlers"
	"lambda-server/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// SetupRouter initializes and configures the Gin router with all routes
// SetupRouter initializes and configures the Gin router with all routes
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Middleware
	r.Use(corsMiddleware())

	// Health check
	r.GET("/health", healthCheck)

	// API routes
	api := r.Group("/api")
	{
		//auth login google logout register
		SetupAuthRoutes(api)

		// Setup user routes
		// SetupUserRoutes(api)

		// Setup journal routes
		SetupJournalRoutes(api)

		// Setup emergency routes
		SetupEmergencyRoutes(api)

		// Setup survey routes
		SetupSurveyRoutes(api)

		// Setup chat routes
		SetupChatRoutes(api)
	}

	return r
}

// corsMiddleware handles CORS headers
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		print("Origin received: ", origin)
		var allowedOrigin string
		if utils.IsRunningLocally() {
			allowedOrigin = "http://localhost:3000"
		} else {
			allowedOrigin = "https://main.d2l1lly6wpq28n.amplifyapp.com"
		}
		c.Header("Access-Control-Allow-Origin", allowedOrigin)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Length, Content-Type, Authorization, X-Requested-With, Accept, Accept-Encoding, Accept-Language, Cache-Control, X-CSRF-Token, X-Client-Type")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Type, X-New-Access-Token, X-New-Refresh-Token")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	}
}

// healthCheck provides a health endpoint
func healthCheck(c *gin.Context) {
	env := "lambda"
	if utils.IsRunningLocally() {
		env = "local"
	}
	c.JSON(http.StatusOK, gin.H{
		"status":      "healthy",
		"environment": env,
		"timestamp":   time.Now().UTC().Format(time.RFC3339),
	})
}

// SetupChatRoutes registers chat-related endpoints
func SetupChatRoutes(rg *gin.RouterGroup) {
	rg.POST("/chat", handlers.HandleChat)
}
