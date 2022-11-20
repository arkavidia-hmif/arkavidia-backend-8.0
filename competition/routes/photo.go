package routes

import (
	"github.com/gin-gonic/gin"

	"arkavidia-backend-8.0/competition/controllers"
)

func PhotoRoute(route *gin.Engine) {
	route.GET("/get-photo-data", controllers.AddPhotoHandler())
	route.POST("/add-photo", controllers.AddPhotoHandler())
	route.DELETE("/delete-photo", controllers.DeletePhotoHandler())
}
