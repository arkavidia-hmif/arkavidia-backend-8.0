package routes

import (
	"github.com/gin-gonic/gin"

	"arkavidia-backend-8.0/competition/controllers"
	"arkavidia-backend-8.0/competition/middlewares"
)

func SubmissionRoute(route *gin.Engine) {
	submissionGroup := route.Group("/submission")

	submissionGroup.GET("/get", middlewares.AuthMiddleware(), controllers.GetSubmissionHandler())
	submissionGroup.GET("/get-all", middlewares.AdminMiddleware(), controllers.GetAllSubmissionsHandler())
	submissionGroup.GET("/download", middlewares.AdminMiddleware(), controllers.DownloadSubmissionHandler())
	submissionGroup.POST("/add", middlewares.AuthMiddleware(), controllers.AddSubmissionHandler())
	submissionGroup.DELETE("/delete", middlewares.AuthMiddleware(), controllers.DeleteSubmissionHandler())
}
