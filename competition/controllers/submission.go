package controllers

import (
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	// storageConfig "arkavidia-backend-8.0/competition/config/storage"
	"arkavidia-backend-8.0/competition/models"
	databaseService "arkavidia-backend-8.0/competition/services/database"
	// storageService "arkavidia-backend-8.0/competition/services/storage"
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
		teamID := c.MustGet("team_id").(uuid.UUID)

		condition := models.Submission{TeamID: teamID}
		submissions := []models.Submission{}
		if err := db.Where(&condition).Find(&submissions).Error; err != nil {
			response := gin.H{"Message": "Error: Bad Request!"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		response := gin.H{"Message": "Success", "Data": submissions} // Ubah tipe data dari response
		c.JSON(http.StatusOK, response)
		return
	}
}

func AddSubmissionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// db := databaseService.GetDB()
		// storage := storageService.GetClient()
		// config := storageConfig.GetStorageConfig()
		// teamID := c.MustGet("team_id").(uuid.UUID)
		// request := AddSubmissionRequest{}
	}
}

func EditSubmissionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// db := databaseService.GetDB()
		// storage := storageService.GetClient()
		// config := storageConfig.GetStorageConfig()
		// teamID := c.MustGet("team_id").(uuid.UUID)
		// request := EditSubmissionRequest{}
	}
}
