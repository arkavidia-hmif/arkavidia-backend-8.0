package controllers

import (
	"mime/multipart"

	"arkavidia-backend-8.0/competition/models"
	"github.com/gin-gonic/gin"
)

type AddSubmissionRequest struct {
	File *multipart.FileHeader `form:"file" field:"file" binding:"required"`
}

type EditSubmissionRequest struct {
	Stage models.SubmissionStage `from:"stage" field:"stage"`
	File  *multipart.FileHeader  `form:"file" field:"file" binding:"required"`
}

func GetSubmissionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func AddSubmissionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func EditSubmissionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
