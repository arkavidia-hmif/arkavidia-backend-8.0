package controllers

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"

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

type EditSubmissionRequest struct {
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

		response := gin.H{"Message": "Success", "Data": submissions, "URL": fmt.Sprintf("https://storage.googleapis.com/arkavidia-8/%s/", config.PhotoDir)}
		c.JSON(http.StatusOK, response)
		return
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

		filename := uuid.New()
		fileext := filepath.Ext(request.File.Filename)

		submission := models.Submission{FileName: filename, FileExtension: fileext, TeamID: teamID, Timestamp: time.Now(), Stage: request.Stage}
		if err := db.Create(&submission).Error; err != nil {
			response := gin.H{"Message": "Error: Bad Request!"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		if err := storageService.UploadFile(client, fmt.Sprintf("%s%s", filename, fileext), config.SubmissionDir, openedFile); err != nil {
			response := gin.H{"Message": "Error: Google Cloud Storage Cannot be Accessed!"}
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		response := gin.H{"Message": "Success", "Data": submission, "URL": fmt.Sprintf("https://storage.googleapis.com/arkavidia-8/%s/", config.PhotoDir)}
		c.JSON(http.StatusCreated, response)
		return
	}
}

func EditSubmissionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// db := databaseService.GetDB()
		// client := storageService.GetClient()
		// config := storageConfig.GetStorageConfig()
		// teamID := c.MustGet("team_id").(uuid.UUID)
		// request := EditSubmissionRequest{}
	}
}
