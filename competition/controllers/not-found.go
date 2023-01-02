package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"arkavidia-backend-8.0/competition/repository"
)

func NotFoundHanlder() gin.HandlerFunc {
	return func(c *gin.Context) {
		response := repository.Response[string]{}
		response.Message = "ERROR: PATH NOT FOUND"
		c.AbortWithStatusJSON(http.StatusNotFound, response)
	}
}
