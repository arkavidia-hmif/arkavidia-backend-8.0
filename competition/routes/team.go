package routes

import (
	"github.com/gin-gonic/gin"

	"arkavidia-backend-8.0/competition/controllers"
	"arkavidia-backend-8.0/competition/middlewares"
)

func TeamRoute(route *gin.Engine) {
	route.POST("/sign-in", controllers.SignInHandler())
	route.POST("/sign-up", controllers.SignUpHandler())
	route.GET("/get-team-data", middlewares.AuthMiddleware(), controllers.GetTeamHandler())
	route.PUT("/change-password", middlewares.AuthMiddleware(), controllers.ChangePasswordHandler())
	route.POST("/competition-registration", middlewares.AuthMiddleware(), controllers.CompetitionRegistration())
}
