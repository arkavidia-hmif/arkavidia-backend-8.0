package routes

import (
	"github.com/gin-gonic/gin"

	"arkavidia-backend-8.0/competition/controllers"
	"arkavidia-backend-8.0/competition/middlewares"
)

func ParticipantRoute(route *gin.Engine) {
	route.GET("/get-member", middlewares.AuthMiddleware(), controllers.GetMemberHandler())
	route.POST("/add-member", middlewares.AuthMiddleware(), controllers.AddMemberHandler())
	route.PUT("/change-career-interest", middlewares.AuthMiddleware(), controllers.ChangeCareerInterestHandler())
	route.PUT("/change-role", middlewares.AuthMiddleware(), controllers.ChangeRoleInterestHandler())
	route.DELETE("/delete-member", middlewares.AuthMiddleware(), controllers.DeleteParticipantHandler())
}
