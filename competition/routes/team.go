package routes

import (
	"github.com/gin-gonic/gin"

	"arkavidia-backend-8.0/competition/controllers"
	"arkavidia-backend-8.0/competition/middlewares"
)

func TeamRoute(route *gin.Engine) {
	groupTeam := route.Group("/team")

	groupTeam.GET("/get", middlewares.AuthMiddleware(), controllers.GetTeamHandler())
	groupTeam.GET("/get-all", middlewares.AdminMiddleware(), controllers.GetAllTeamsHandler())
	groupTeam.POST("/sign-in", controllers.SignInTeamHandler())
	groupTeam.POST("/sign-up", controllers.SignUpTeamHandler())
	groupTeam.PUT("/change-password", middlewares.AuthMiddleware(), controllers.ChangePasswordHandler())
	groupTeam.PUT("/competition-registration", middlewares.AuthMiddleware(), controllers.CompetitionRegistration())
	groupTeam.PUT("/change-status", middlewares.AdminMiddleware(), controllers.ChangeStatusTeamHandler())
}
