package routes

import (
	"github.com/gin-gonic/gin"

	"arkavidia-backend-8.0/competition/controllers"
)

func NotFoundRoute(route *gin.Engine) {
	route.NoRoute(controllers.NotFoundHanlder())
}
