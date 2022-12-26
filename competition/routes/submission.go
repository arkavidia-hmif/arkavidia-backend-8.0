package routes

import (
	"github.com/gin-gonic/gin"

	"arkavidia-backend-8.0/competition/controllers"
	"arkavidia-backend-8.0/competition/middlewares"
	"arkavidia-backend-8.0/competition/utils/cache"
)

func SubmissionRoute(route *gin.Engine) {
	submissionGroup := route.Group("/submission")

	submissionGroup.GET("/", middlewares.AuthMiddleware(), cache.Store.GetHandlerFunc(controllers.GetSubmissionHandler()))
	submissionGroup.GET("/all", middlewares.AuthMiddleware(), cache.Store.GetHandlerFunc(controllers.GetAllSubmissionsHandler()))
	submissionGroup.GET("/download", middlewares.AuthMiddleware(), cache.Store.GetHandlerFunc(controllers.DownloadSubmissionHandler()))
	submissionGroup.POST("/", middlewares.AuthMiddleware(), controllers.AddSubmissionHandler())
	submissionGroup.DELETE("/", middlewares.AuthMiddleware(), controllers.DeleteSubmissionHandler())
}
