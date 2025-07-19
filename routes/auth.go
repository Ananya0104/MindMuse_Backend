package routes

import (
	"lambda-server/handlers"
	"lambda-server/middlewares"

	"github.com/gin-gonic/gin"
)

// SetupUserRoutes configures all user-related routes
func SetupAuthRoutes(api *gin.RouterGroup) {
	user := api.Group("/auth")
	{
		user.POST("/login", handlers.HandleLogin)
		user.POST("/register", handlers.HandleRegister)
		user.POST("/refresh", handlers.HandleRefresh)
		user.POST("/logout", middlewares.AuthMiddleware(), handlers.HandleLogout)
		user.GET("/me", handlers.HandleGetProfile)
		user.PATCH("/me", handlers.UpdateCurrentUser)
		user.DELETE("/me", handlers.DeleteCurrentUser)
		user.POST("/forgot-password", handlers.HandleForgotPassword)
		user.POST("/reset-password", handlers.HandleResetPassword)
	}
}
