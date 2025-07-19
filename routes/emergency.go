package routes

import (
	"lambda-server/handlers"
	"lambda-server/middlewares"

	"github.com/gin-gonic/gin"
)

// SetupUserRoutes configures all user-related routes
func SetupEmergencyRoutes(api *gin.RouterGroup) {
	emergency := api.Group("/emergency")
	emergency.POST("/create", middlewares.AuthMiddleware(), handlers.CreateEmergencyContacts)
	emergency.GET("/contacts", middlewares.AuthMiddleware(), handlers.GetEmergencyContacts)
}
