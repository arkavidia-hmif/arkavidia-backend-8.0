package routes

import (
	"github.com/gin-gonic/gin"

	"arkavidia-backend-8.0/competition/controllers"
	"arkavidia-backend-8.0/competition/middlewares"
)

func SubmissionRoute(route *gin.Engine) {
	route.GET("/get-submission", middlewares.AuthMiddleware(), controllers.GetSubmissionHandler())
	route.POST("/add-submission", middlewares.AuthMiddleware(), controllers.AddSubmissionHandler())
	route.DELETE("/delete-submission", middlewares.AuthMiddleware(), controllers.DeleteSubmissionHandler())
}
