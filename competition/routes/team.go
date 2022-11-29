package routes

import (
	"github.com/gin-gonic/gin"

	"arkavidia-backend-8.0/competition/controllers"
	"arkavidia-backend-8.0/competition/middlewares"
)

func TeamRoute(route *gin.Engine) {
	route.GET("/get-team", middlewares.AuthMiddleware(), controllers.GetTeamHandler())
	route.POST("/sign-in", controllers.SignInHandler())
	route.POST("/sign-up", controllers.SignUpHandler())
	route.PUT("/change-password", middlewares.AuthMiddleware(), controllers.ChangePasswordHandler())
	route.PUT("/competition-registration", middlewares.AuthMiddleware(), controllers.CompetitionRegistration())
}
