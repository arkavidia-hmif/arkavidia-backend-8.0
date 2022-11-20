package routes

import (
	"github.com/gin-gonic/gin"

	"arkavidia-backend-8.0/competition/controllers"
)

func ParticipantRoute(route *gin.Engine) {
	route.GET("/get-participant-data", controllers.GetParticipantHandler())
	route.POST("/add-participant", controllers.AddParticipantHandler())
}
