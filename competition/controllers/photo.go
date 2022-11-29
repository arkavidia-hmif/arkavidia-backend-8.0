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
	ParticipantID uint `json:"participant_id" binding:"required"`
}

type AddPhotoRequest struct {
	ParticipantID uint                  `form:"participant_id" field:"participant_id" binding:"required"`
	Type          models.PhotoType      `form:"type" field:"type" binding:"required"`
	File          *multipart.FileHeader `form:"file" field:"file" binding:"required"`
}

type DeletePhotoRequest struct {
	ParticipantID uint   `json:"participant_id" binding:"required"`
	FileName      string `json:"file_name" binding:"required"`
	FileExtension string `json:"file_extension" binding:"required"`
}

func GetPhotoHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.GetDB()
		config := storageConfig.GetStorageConfig()
		teamID := c.MustGet("team_id").(uint)

		request := GetPhotoRequest{}
		if err := c.BindJSON(&request); err != nil {
			response := gin.H{"Message": "ERROR: BAD REQUEST"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		condition1 := models.Membership{TeamID: teamID, ParticipantID: request.ParticipantID}
		membership := models.Membership{TeamID: teamID, ParticipantID: request.ParticipantID}
		if err := db.Where(&condition1).Find(&membership).Error; err != nil {
			response := gin.H{"Message": "ERROR: BAD REQUEST"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		condition2 := models.Photo{ParticipantID: request.ParticipantID}
		photos := []models.Photo{}
		if err := db.Where(&condition2).Find(&photos).Error; err != nil {
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
		teamID := c.MustGet("team_id").(uint)

		request := AddPhotoRequest{}
		if err := c.MustBindWith(&request, binding.FormMultipart); err != nil {
			response := gin.H{"Message": "Error: Bad Request"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		condition := models.Membership{TeamID: teamID, ParticipantID: request.ParticipantID}
		membership := models.Membership{}
		if err := db.Where(&condition).Find(&membership).Error; err != nil {
			response := gin.H{"Message": "ERROR: BAD REQUEST"}
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

		if err := storageService.UploadFile(client, fmt.Sprintf("%s%s", fileUUID, fileExt), config.PhotoDir, openedFile); err != nil {
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
		teamID := c.MustGet("team_id").(uint)

		request := DeletePhotoRequest{}
		if err := c.BindJSON(&request); err != nil {
			response := gin.H{"Message": "ERROR: BAD REQUEST"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		condition1 := models.Membership{TeamID: teamID, ParticipantID: request.ParticipantID}
		membership := models.Membership{}
		if err := db.Where(&condition1).Find(&membership).Error; err != nil {
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
		condition2 := models.Photo{FileName: fileUUID}
		photo := models.Photo{}
		if err := db.Where(&condition2).Delete(&photo).Error; err != nil {
			response := gin.H{"Message": "ERROR: BAD REQUEST"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		response := gin.H{"Message": "SUCCESS"}
		c.JSON(http.StatusOK, response)
	}
}
