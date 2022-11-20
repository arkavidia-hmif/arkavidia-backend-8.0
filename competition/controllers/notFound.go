package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NotFoundHanlder() gin.HandlerFunc {
	return func(c *gin.Context) {
		response := gin.H{"Message": "ERROR: PATH NOT FOUND"}
		c.JSON(http.StatusNotFound, response)
	}
}
