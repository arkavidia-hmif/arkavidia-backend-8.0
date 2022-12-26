package routes

import (
	"github.com/gin-gonic/gin"

	"arkavidia-backend-8.0/competition/controllers"
	"arkavidia-backend-8.0/competition/middlewares"
	"arkavidia-backend-8.0/competition/utils/cache"
)

func TeamRoute(route *gin.Engine) {
	groupTeam := route.Group("/team")

	groupTeam.GET("/", middlewares.AuthMiddleware(), cache.Store.GetHandlerFunc(controllers.GetTeamHandler()))
	groupTeam.GET("/all", middlewares.AuthMiddleware(), cache.Store.GetHandlerFunc(controllers.GetAllTeamsHandler()))
	groupTeam.POST("/sign-in", controllers.SignInTeamHandler())
	groupTeam.POST("/", controllers.SignUpTeamHandler())
	groupTeam.PUT("/password", middlewares.AuthMiddleware(), controllers.ChangePasswordHandler())
	groupTeam.PUT("/registration", middlewares.AuthMiddleware(), controllers.CompetitionRegistration())
	groupTeam.PUT("/status", middlewares.AuthMiddleware(), controllers.ChangeStatusTeamHandler())
}
