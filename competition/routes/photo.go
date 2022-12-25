package routes

import (
	"github.com/gin-gonic/gin"

	"arkavidia-backend-8.0/competition/controllers"
	"arkavidia-backend-8.0/competition/middlewares"
)

func PhotoRoute(route *gin.Engine) {
	photoGroup := route.Group("/photo")

	photoGroup.GET("/", middlewares.AuthMiddleware(), controllers.GetPhotoHandler())
	photoGroup.GET("/all", middlewares.AuthMiddleware(), controllers.GetAllPhotosHandler())
	photoGroup.GET("/download", middlewares.AuthMiddleware(), controllers.DownloadPhotoHandler())
	photoGroup.POST("/", middlewares.AuthMiddleware(), controllers.AddPhotoHandler())
	photoGroup.PUT("/status", middlewares.AuthMiddleware(), controllers.ChangeStatusPhotoHandler())
	photoGroup.DELETE("/", middlewares.AuthMiddleware(), controllers.DeletePhotoHandler())
}
