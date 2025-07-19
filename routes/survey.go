package routes

import (
	"lambda-server/handlers"
	"lambda-server/middlewares"

	"github.com/gin-gonic/gin"
)

// SetupSurveyRoutes configures all survey-related routes
func SetupSurveyRoutes(api *gin.RouterGroup) {

	score := api.Group("/score")
	{
		score.POST("/submit", middlewares.AuthMiddleware(), handlers.SubmitMindMuseScore)
	}
}
