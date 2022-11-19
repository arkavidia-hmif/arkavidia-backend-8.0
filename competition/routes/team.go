package routes

import (
	"github.com/gin-gonic/gin"

	"arkavidia-backend-8.0/competition/controllers"
)

func TeamRoute(route *gin.Engine) {
	route.POST("/sign-in", controllers.SignInHandler())
	route.POST("/sign-up", controllers.SignUpHandler())
	route.GET("/get-team-data", controllers.GetTeam())
	route.PUT("/change-password", controllers.ChangePasswordHandler())
	route.POST("/competition-registration", controllers.CompetitionRegistration())
}
