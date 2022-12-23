package routes

import (
	"github.com/gin-gonic/gin"

	"arkavidia-backend-8.0/competition/controllers"
	"arkavidia-backend-8.0/competition/middlewares"
)

func TeamRoute(route *gin.Engine) {
	groupTeam := route.Group("/team")

	groupTeam.GET("/", middlewares.AuthMiddleware(), controllers.GetTeamHandler())
	groupTeam.GET("/all", middlewares.AdminMiddleware(), controllers.GetAllTeamsHandler())
	groupTeam.POST("/sign-in", controllers.SignInTeamHandler())
	groupTeam.POST("/", controllers.SignUpTeamHandler())
	groupTeam.PUT("/password", middlewares.AuthMiddleware(), controllers.ChangePasswordHandler())
	groupTeam.PUT("/registration", middlewares.AuthMiddleware(), controllers.CompetitionRegistration())
	groupTeam.PUT("/status", middlewares.AdminMiddleware(), controllers.ChangeStatusTeamHandler())
}
