package routes

import (
	"github.com/gin-gonic/gin"

	"arkavidia-backend-8.0/competition/controllers"
)

func TeamRoute(route *gin.Engine) {
	route.POST("/sign-in", controllers.SignInHandler())
	route.POST("/sign-up", controllers.SignUpHandler())
	route.GET("/team", controllers.GetTeam())
}
