package controllers

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"

	storageConfig "arkavidia-backend-8.0/competition/config/storage"
	"arkavidia-backend-8.0/competition/models"
	databaseService "arkavidia-backend-8.0/competition/services/database"
	storageService "arkavidia-backend-8.0/competition/services/storage"
)

type AddSubmissionRequest struct {
	Stage models.SubmissionStage `from:"stage" field:"stage"`
	File  *multipart.FileHeader  `form:"file" field:"file" binding:"required"`
}

type DeleteSubmissionRequest struct {
	Stage models.SubmissionStage `from:"stage" field:"stage"`
	File  *multipart.FileHeader  `form:"file" field:"file" binding:"required"`
}

func GetSubmissionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.GetDB()
		config := storageConfig.GetStorageConfig()
		teamID := c.MustGet("team_id").(uuid.UUID)

		condition := models.Submission{TeamID: teamID}
		submissions := []models.Submission{}
		if err := db.Where(&condition).Find(&submissions).Error; err != nil {
			response := gin.H{"Message": "Error: Bad Request!"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		response := gin.H{"Message": "Success", "Data": submissions, "URL": fmt.Sprintf("%s/%s/%s/", config.StorageHost, config.BucketName, config.SubmissionDir)}
		c.JSON(http.StatusOK, response)
	}
}

func AddSubmissionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.GetDB()
		client := storageService.GetClient()
		config := storageConfig.GetStorageConfig()
		teamID := c.MustGet("team_id").(uuid.UUID)
		request := AddSubmissionRequest{}

		if err := c.MustBindWith(&request, binding.FormMultipart); err != nil {
			response := gin.H{"Message": "Error: Bad Request!"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		openedFile, err := request.File.Open()
		if err != nil {
			response := gin.H{"Message": "Error: File Cannot be Accessed!"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		fileUUID := uuid.New()
		fileExt := filepath.Ext(request.File.Filename)

		submission := models.Submission{FileName: fileUUID, FileExtension: fileExt, TeamID: teamID, Stage: request.Stage}
		if err := db.Create(&submission).Error; err != nil {
			response := gin.H{"Message": "Error: Bad Request!"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		if err := storageService.UploadFile(client, fmt.Sprintf("%s%s", fileUUID, fileExt), config.SubmissionDir, openedFile); err != nil {
			response := gin.H{"Message": "Error: Google Cloud Storage Cannot be Accessed!"}
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		response := gin.H{"Message": "Success", "Data": submission, "URL": fmt.Sprintf("%s/%s/%s/", config.StorageHost, config.BucketName, config.SubmissionDir)}
		c.JSON(http.StatusCreated, response)
	}
}

func DeleteSubmissionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.GetDB()
		client := storageService.GetClient()
		config := storageConfig.GetStorageConfig()
		teamID := c.MustGet("team_id").(uuid.UUID)
		request := DeleteSubmissionRequest{}

		if err := c.MustBindWith(&request, binding.FormMultipart); err != nil {
			response := gin.H{"Message": "Error: Bad Request!"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		submission := models.Submission{TeamID: teamID, Stage: request.Stage}
		if err := db.Find(&submission).Error; err != nil {
			response := gin.H{"Message": "Error: Bad Request!"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		if err := storageService.DeleteFile(client, fmt.Sprintf("%s%s", submission.FileName, submission.FileExtension), config.PhotoDir); err != nil {
			response := gin.H{"Message": "Error: Google Cloud Storage Cannot be Accessed!"}
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		response := gin.H{"Message": "Success"}
		c.JSON(http.StatusOK, response)
	}
}
