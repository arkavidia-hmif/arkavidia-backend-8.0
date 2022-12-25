package routes

import (
	"github.com/gin-gonic/gin"

	"arkavidia-backend-8.0/competition/controllers"
	"arkavidia-backend-8.0/competition/middlewares"
)

func SubmissionRoute(route *gin.Engine) {
	submissionGroup := route.Group("/submission")

	submissionGroup.GET("/", middlewares.AuthMiddleware(), controllers.GetSubmissionHandler())
	submissionGroup.GET("/all", middlewares.AuthMiddleware(), controllers.GetAllSubmissionsHandler())
	submissionGroup.GET("/download", middlewares.AuthMiddleware(), controllers.DownloadSubmissionHandler())
	submissionGroup.POST("/", middlewares.AuthMiddleware(), controllers.AddSubmissionHandler())
	submissionGroup.DELETE("/", middlewares.AuthMiddleware(), controllers.DeleteSubmissionHandler())
}
