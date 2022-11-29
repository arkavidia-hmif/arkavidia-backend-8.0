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
	Stage models.SubmissionStage `form:"stage" field:"stage" binding:"required"`
	File  *multipart.FileHeader  `form:"file" field:"file" binding:"required"`
}

type DeleteSubmissionRequest struct {
	FileName      string `json:"file_name" binding:"required"`
	FileExtension string `json:"file_extension" binding:"required"`
}

func GetSubmissionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.GetDB()
		config := storageConfig.GetStorageConfig()
		teamID := c.MustGet("team_id").(uint)

		condition := models.Submission{TeamID: teamID}
		submissions := []models.Submission{}
		if err := db.Where(&condition).Find(&submissions).Error; err != nil {
			response := gin.H{"Message": "ERROR: BAD REQUEST"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		response := gin.H{"Message": "SUCCESS", "Data": submissions, "URL": fmt.Sprintf("%s/%s/%s/", config.StorageHost, config.BucketName, config.SubmissionDir)}
		c.JSON(http.StatusOK, response)
	}
}

func AddSubmissionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.GetDB()
		client := storageService.GetClient()
		config := storageConfig.GetStorageConfig()
		teamID := c.MustGet("team_id").(uint)

		request := AddSubmissionRequest{}
		if err := c.MustBindWith(&request, binding.FormMultipart); err != nil {
			response := gin.H{"Message": "ERROR: BAD REQUEST"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		openedFile, err := request.File.Open()
		if err != nil {
			response := gin.H{"Message": "ERROR: FILE CANNOT BE ACCESSED"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		fileUUID := uuid.New()
		fileExt := filepath.Ext(request.File.Filename)

		submission := models.Submission{FileName: fileUUID, FileExtension: fileExt, TeamID: teamID, Stage: request.Stage}
		if err := db.Create(&submission).Error; err != nil {
			response := gin.H{"Message": "ERROR: BAD REQUEST"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		if err := storageService.UploadFile(client, fmt.Sprintf("%s%s", fileUUID, fileExt), config.SubmissionDir, openedFile); err != nil {
			response := gin.H{"Message": "ERROR: GOOGLE CLOUD STORAGE CANNOT BE ACCESSED"}
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		response := gin.H{"Message": "SUCCESS", "Data": submission, "URL": fmt.Sprintf("%s/%s/%s/", config.StorageHost, config.BucketName, config.SubmissionDir)}
		c.JSON(http.StatusCreated, response)
	}
}

func DeleteSubmissionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.GetDB()
		teamID := c.MustGet("team_id").(uint)

		request := DeleteSubmissionRequest{}
		if err := c.BindJSON(&request); err != nil {
			response := gin.H{"Message": "ERROR: BAD REQUEST"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		fileUUID, err := uuid.Parse(request.FileName)
		if err != nil {
			response := gin.H{"Message": "ERROR: BAD REQUEST"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		condition := models.Submission{FileName: fileUUID, TeamID: teamID}
		submission := models.Submission{}
		if err := db.Where(&condition).Delete(&submission).Error; err != nil {
			response := gin.H{"Message": "ERROR: BAD REQUEST"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		response := gin.H{"Message": "SUCCESS"}
		c.JSON(http.StatusOK, response)
	}
}
