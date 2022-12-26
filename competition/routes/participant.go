package routes

import (
	"github.com/gin-gonic/gin"

	"arkavidia-backend-8.0/competition/controllers"
	"arkavidia-backend-8.0/competition/middlewares"
	"arkavidia-backend-8.0/competition/utils/cache"
)

func ParticipantRoute(route *gin.Engine) {
	participantGroup := route.Group("/participant")

	participantGroup.GET("/", middlewares.AuthMiddleware(), cache.Store.GetHandlerFunc(controllers.GetMemberHandler()))
	participantGroup.GET("/all", middlewares.AuthMiddleware(), cache.Store.GetHandlerFunc(controllers.GetAllMembersHandler()))
	participantGroup.POST("/", middlewares.AuthMiddleware(), controllers.AddMemberHandler())
	participantGroup.PUT("/career-interest", middlewares.AuthMiddleware(), controllers.ChangeCareerInterestHandler())
	participantGroup.PUT("/role", middlewares.AuthMiddleware(), controllers.ChangeRoleHandler())
	participantGroup.PUT("/status", middlewares.AuthMiddleware(), controllers.ChangeStatusParticipantHandler())
	participantGroup.DELETE("/", middlewares.AuthMiddleware(), controllers.DeleteParticipantHandler())
}
