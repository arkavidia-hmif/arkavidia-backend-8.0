package routes

import (
	"github.com/gin-gonic/gin"

	"arkavidia-backend-8.0/competition/controllers"
	"arkavidia-backend-8.0/competition/middlewares"
)

func PhotoRoute(route *gin.Engine) {
	photoGroup := route.Group("/photo")

	photoGroup.GET("/get", middlewares.AuthMiddleware(), controllers.GetPhotoHandler())
	photoGroup.GET("/get-all", middlewares.AdminMiddleware(), controllers.GetAllPhotosHandler())
	photoGroup.GET("/download", middlewares.AdminMiddleware(), controllers.DownloadPhotoHandler())
	photoGroup.POST("/add", middlewares.AuthMiddleware(), controllers.AddPhotoHandler())
	photoGroup.PUT("/change-status", middlewares.AdminMiddleware(), controllers.ChangeStatusPhotoHandler())
	photoGroup.DELETE("/delete", middlewares.AuthMiddleware(), controllers.DeletePhotoHandler())
}
