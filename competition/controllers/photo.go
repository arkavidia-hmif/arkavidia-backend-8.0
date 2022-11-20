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

type GetPhotoRequest struct {
	ParticipantID uint `json:"participant_id"`
}

type AddPhotoRequest struct {
	ParticipantID uint                  `form:"participant_id" field:"participant_id"`
	Type          models.PhotoType      `form:"type" field:"type"`
	File          *multipart.FileHeader `form:"file" field:"file" binding:"required"`
}

type DeletePhotoRequest struct {
	ParticipantID uint                  `form:"participant_id" field:"participant_id"`
	Type          models.PhotoType      `form:"type" field:"type"`
	File          *multipart.FileHeader `form:"file" field:"file" binding:"required"`
}

func GetPhotoHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.GetDB()
		config := storageConfig.GetStorageConfig()

		request := GetPhotoRequest{}
		if err := c.BindJSON(&request); err != nil {
			response := gin.H{"Message": "Error: Bad Request"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		condition := models.Photo{ParticipantID: request.ParticipantID}
		photos := []models.Photo{}
		if err := db.Where(&condition).Find(&photos).Error; err != nil {
			response := gin.H{"Message": "Error: Bad Request"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		response := gin.H{"Message": "Success", "Data": photos, "URL": fmt.Sprintf("%s/%s/%s", config.StorageHost, config.BucketName, config.PhotoDir)}
		c.JSON(http.StatusOK, response)
	}
}

func AddPhotoHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.GetDB()
		client := storageService.GetClient()
		config := storageConfig.GetStorageConfig()

		request := AddPhotoRequest{}
		if err := c.MustBindWith(&request, binding.FormMultipart); err != nil {
			response := gin.H{"Message": "Error: Bad Request"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		openedFile, err := request.File.Open()
		if err != nil {
			response := gin.H{"Message": "Error: File Cannot be Accessed"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		fileUUID := uuid.New()
		fileExt := filepath.Ext(request.File.Filename)

		photo := models.Photo{FileName: fileUUID, FileExtension: fileExt, ParticipantID: request.ParticipantID, Status: models.WaitingForVerification, Type: request.Type}
		if err := db.Create(&photo).Error; err != nil {
			response := gin.H{"Message": "Error: Bad Request"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		if err := storageService.UploadFile(client, fmt.Sprintf("%s%s", fileUUID, fileExt), config.SubmissionDir, openedFile); err != nil {
			response := gin.H{"Message": "Error: Google Cloud Storage Cannot be Accessed"}
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		response := gin.H{"Message": "Success", "Data": photo, "URL": fmt.Sprintf("%s/%s/%s", config.StorageHost, config.BucketName, config.PhotoDir)}
		c.JSON(http.StatusCreated, response)
	}
}

func DeletePhotoHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.GetDB()
		client := storageService.GetClient()
		config := storageConfig.GetStorageConfig()

		request := DeletePhotoRequest{}
		if err := c.MustBindWith(&request, binding.FormMultipart); err != nil {
			response := gin.H{"Message": "ERROR: BAD REQUEST"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		photo := models.Photo{ParticipantID: request.ParticipantID, Type: request.Type}
		if err := db.Create(&photo).Error; err != nil {
			response := gin.H{"Message": "ERROR: BAD REQUEST"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		if err := storageService.DeleteFile(client, fmt.Sprintf("%s%s", photo.FileName, photo.FileExtension), config.PhotoDir); err != nil {
			response := gin.H{"Message": "ERROR: GOOGLE CLOUD STORAGE CANNOT BE ACCESSED"}
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		response := gin.H{"Message": "SUCCESS"}
		c.JSON(http.StatusOK, response)
	}
}
