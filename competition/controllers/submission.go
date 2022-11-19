package controllers

import (
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	storageConfig "arkavidia-backend-8.0/competition/config/storage"
	"arkavidia-backend-8.0/competition/models"
	databaseService "arkavidia-backend-8.0/competition/services/database"
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
		db := databaseService.GetDB()
		config := storageConfig.GetStorageConfig()
		teamID := c.MustGet("team_id").(uuid.UUID)
	}
}

func AddSubmissionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.GetDB()
		config := storageConfig.GetStorageConfig()
		teamID := c.MustGet("team_id").(uuid.UUID)
		request := AddSubmissionRequest{}
	}
}

func EditSubmissionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.GetDB()
		config := storageConfig.GetStorageConfig()
		teamID := c.MustGet("team_id").(uuid.UUID)
		request := EditSubmissionRequest{}
	}
}
