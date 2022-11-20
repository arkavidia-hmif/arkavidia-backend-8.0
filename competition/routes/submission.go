package routes

import (
	"github.com/gin-gonic/gin"

	"arkavidia-backend-8.0/competition/controllers"
)

func SubmissionRoute(route *gin.Engine) {
	route.GET("/get-submission-data", controllers.GetSubmissionHandler())
	route.POST("/add-submission", controllers.AddSubmissionHandler())
	route.DELETE("/delete-submission", controllers.DeleteSubmissionHandler())
}
