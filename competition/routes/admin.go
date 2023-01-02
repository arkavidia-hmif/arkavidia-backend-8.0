package routes

import (
	"github.com/gin-gonic/gin"

	"arkavidia-backend-8.0/competition/controllers"
)

func AdminRoute(route *gin.Engine) {
	adminGroup := route.Group("/admin")

	// NOTE: Penambahan admin baru dilakukan langsung pada basis data
	adminGroup.POST("/sign-in", controllers.SignInAdminHandler())
}
