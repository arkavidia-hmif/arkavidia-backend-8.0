package routes

import (
	"github.com/gin-gonic/gin"

	"arkavidia-backend-8.0/competition/controllers"
	"arkavidia-backend-8.0/competition/middlewares"
)

func PhotoRoute(route *gin.Engine) {
	route.GET("/get-photo", middlewares.AuthMiddleware(), controllers.GetPhotoHandler())
	route.POST("/get-photo", middlewares.AuthMiddleware(), controllers.GetPhotoHandler())
	route.POST("/add-photo", middlewares.AuthMiddleware(), controllers.AddPhotoHandler())
	route.DELETE("/delete-photo", middlewares.AuthMiddleware(), controllers.DeletePhotoHandler())
}
