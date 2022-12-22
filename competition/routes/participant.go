package routes

import (
	"github.com/gin-gonic/gin"

	"arkavidia-backend-8.0/competition/controllers"
	"arkavidia-backend-8.0/competition/middlewares"
)

func ParticipantRoute(route *gin.Engine) {
	participantGroup := route.Group("/member")

	participantGroup.GET("/get", middlewares.AuthMiddleware(), controllers.GetMemberHandler())
	participantGroup.GET("/get-all", middlewares.AdminMiddleware(), controllers.GetAllMembersHandler())
	participantGroup.POST("/add", middlewares.AuthMiddleware(), controllers.AddMemberHandler())
	participantGroup.PUT("/change-career-interest", middlewares.AuthMiddleware(), controllers.ChangeCareerInterestHandler())
	participantGroup.PUT("/change-role", middlewares.AuthMiddleware(), controllers.ChangeRoleHandler())
	participantGroup.PUT("/change-status", middlewares.AuthMiddleware(), controllers.ChangeStatusParticipantHandler())
	participantGroup.DELETE("/delete", middlewares.AuthMiddleware(), controllers.DeleteParticipantHandler())
}
