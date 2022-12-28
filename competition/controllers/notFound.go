package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"arkavidia-backend-8.0/competition/utils/sanitizer"
)

func NotFoundHanlder() gin.HandlerFunc {
	return func(c *gin.Context) {
		response := sanitizer.Response[string]{}
		response.Message = "ERROR: PATH NOT FOUND"
		c.AbortWithStatusJSON(http.StatusNotFound, sanitizer.SanitizeStruct(response))
	}
}
