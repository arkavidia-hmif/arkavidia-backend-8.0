package routes

import (
	"github.com/gin-gonic/gin"

	"arkavidia-backend-8.0/competition/controllers"
)

func Memebership(route *gin.Engine) {
	route.GET("/get-membership", controllers.GetMember())
	route.POST("/add-membership", controllers.AddMembership())
}
