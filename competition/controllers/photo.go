package controllers

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"gorm.io/gorm"

	storageConfig "arkavidia-backend-8.0/competition/config/storage"
	"arkavidia-backend-8.0/competition/middlewares"
	"arkavidia-backend-8.0/competition/models"
	databaseService "arkavidia-backend-8.0/competition/services/database"
	storageService "arkavidia-backend-8.0/competition/services/storage"
)

type GetPhotoQuery struct {
	ParticipantID uint `form:"participant_id" field:"participant_id" binding:"required"`
}

type GetAllPhotosQuery struct {
	Page int `form:"page" field:"page" binding:"required"`
	Size int `form:"size" field:"size" binding:"required"`
}

type DownloadPhotoQuery struct {
	PhotoID uint `form:"photo_id" field:"photo_id" binding:"required"`
}

type AddPhotoRequest struct {
	ParticipantID uint                  `form:"participant_id" field:"participant_id" binding:"required"`
	Type          models.PhotoType      `form:"type" field:"type" binding:"required"`
	File          *multipart.FileHeader `form:"file" field:"file" binding:"required"`
}

type ChangeStatusPhotoQuery struct {
	PhotoID uint `form:"photo_id" field:"photo_id" binding:"required"`
}

type ChangeStatusPhotoRequest struct {
	Status models.PhotoStatus `json:"status" binding:"required"`
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
		role := c.MustGet("role").(middlewares.AuthRole)

		switch role {
		case middlewares.Admin:
			{
				query := GetPhotoQuery{}
				if err := c.BindQuery(&query); err != nil {
					response := gin.H{"Message": "ERROR: BAD REQUEST"}
					c.JSON(http.StatusBadRequest, response)
					return
				}

				conditionPhoto := models.Photo{ParticipantID: query.ParticipantID}
				photos := []models.Photo{}
				if err := db.Where(&conditionPhoto).Find(&photos).Error; err != nil {
					response := gin.H{"Message": "ERROR: BAD REQUEST"}
					c.JSON(http.StatusBadRequest, response)
					return
				}

				response := gin.H{"Message": "Success", "Data": photos, "URL": fmt.Sprintf("%s/%s/%s", config.StorageHost, config.BucketName, config.PhotoDir)}
				c.JSON(http.StatusOK, response)
				return
			}
		case middlewares.Team:
			{
				query := GetPhotoQuery{}
				if err := c.BindQuery(&query); err != nil {
					response := gin.H{"Message": "ERROR: BAD REQUEST"}
					c.JSON(http.StatusBadRequest, response)
					return
				}

				teamID := c.MustGet("id").(uint)
				conditionMembership := models.Membership{TeamID: teamID, ParticipantID: query.ParticipantID}
				membership := models.Membership{TeamID: teamID, ParticipantID: query.ParticipantID}
				if err := db.Where(&conditionMembership).Find(&membership).Error; err != nil {
					response := gin.H{"Message": "ERROR: BAD REQUEST"}
					c.JSON(http.StatusBadRequest, response)
					return
				}

				conditionPhoto := models.Photo{ParticipantID: query.ParticipantID}
				photos := []models.Photo{}
				if err := db.Where(&conditionPhoto).Find(&photos).Error; err != nil {
					response := gin.H{"Message": "ERROR: BAD REQUEST"}
					c.JSON(http.StatusBadRequest, response)
					return
				}

				response := gin.H{"Message": "Success", "Data": photos, "URL": fmt.Sprintf("%s/%s/%s", config.StorageHost, config.BucketName, config.PhotoDir)}
				c.JSON(http.StatusOK, response)
				return
			}
		default:
			{
				response := gin.H{"Message": "ERROR: INVALID ROLE"}
				c.JSON(http.StatusUnauthorized, response)
				return
			}
		}
	}
}

func GetAllPhotosHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.GetDB()
		role := c.MustGet("role").(middlewares.AuthRole)

		switch role {
		case middlewares.Admin:
			{
				query := GetAllPhotosQuery{}
				if err := c.BindQuery(&query); err != nil {
					response := gin.H{"Message": "ERROR: BAD REQUEST"}
					c.JSON(http.StatusBadRequest, response)
					return
				}

				offset := (query.Page - 1) * query.Size
				limit := query.Size
				photos := []models.Photo{}

				if err := db.Offset(offset).Limit(limit).Find(&photos).Error; err != nil {
					response := gin.H{"Message": "ERROR: BAD REQUEST"}
					c.JSON(http.StatusBadRequest, response)
					return
				}

				response := gin.H{"Message": "SUCCESS", "Data": photos}
				c.JSON(http.StatusOK, response)
				return
			}
		default:
			{
				response := gin.H{"Message": "ERROR: INVALID ROLE"}
				c.JSON(http.StatusUnauthorized, response)
				return
			}
		}
	}
}

func DownloadPhotoHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.GetDB()
		client := storageService.GetClient()
		config := storageConfig.GetStorageConfig()
		role := c.MustGet("role").(middlewares.AuthRole)

		switch role {
		case middlewares.Admin:
			{
				query := DownloadPhotoQuery{}
				if err := c.BindQuery(&query); err != nil {
					response := gin.H{"Message": "ERROR: BAD REQUEST"}
					c.JSON(http.StatusBadRequest, response)
					return
				}

				conditionPhoto := models.Photo{Model: gorm.Model{ID: query.PhotoID}}
				photo := models.Photo{}
				if err := db.Where(&conditionPhoto).Find(&photo).Error; err != nil {
					response := gin.H{"Message": "ERROR: BAD REQUEST"}
					c.JSON(http.StatusBadRequest, response)
					return
				}

				filename := fmt.Sprintf("%s.%s", photo.FileName, photo.FileExtension)
				IOWriter, err := storageService.DownloadFile(client, filename, config.PhotoDir)
				if err != nil {
					response := gin.H{"Message": "ERROR: BAD REQUEST"}
					c.JSON(http.StatusBadRequest, response)
					return
				}

				var content []byte
				length, err := IOWriter.Write(content)
				if err != nil {
					response := gin.H{"Message": "ERROR: INTERNAL SERVER ERROR"}
					c.JSON(http.StatusInternalServerError, response)
					return
				}

				c.Header("Content-Description", "File Transfer")
				c.Header("Content-Transfer-Encoding", "binary")
				c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
				c.Header("Content-Type", "application/octet-stream")
				c.Header("Accept-Length", fmt.Sprintf("%d", length))
				c.Writer.Write(content)

				response := gin.H{"Message": "SUCCESS"}
				c.JSON(http.StatusOK, response)
				return
			}
		default:
			{
				response := gin.H{"Message": "ERROR: INVALID ROLE"}
				c.JSON(http.StatusUnauthorized, response)
				return
			}
		}
	}
}

func AddPhotoHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.GetDB()
		client := storageService.GetClient()
		config := storageConfig.GetStorageConfig()
		role := c.MustGet("role").(middlewares.AuthRole)

		switch role {
		case middlewares.Team:
			{
				request := AddPhotoRequest{}
				if err := c.MustBindWith(&request, binding.FormMultipart); err != nil {
					response := gin.H{"Message": "Error: Bad Request"}
					c.JSON(http.StatusBadRequest, response)
					return
				}

				teamID := c.MustGet("id").(uint)
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
				defer openedFile.Close()

				fileUUID := uuid.New()
				fileExt := filepath.Ext(request.File.Filename)

				photo := models.Photo{FileName: fileUUID, FileExtension: fileExt, ParticipantID: request.ParticipantID, Status: models.WaitingForApproval, Type: request.Type}
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
				return
			}
		default:
			{
				response := gin.H{"Message": "ERROR: INVALID ROLE"}
				c.JSON(http.StatusUnauthorized, response)
				return
			}
		}
	}
}

func ChangeStatusPhotoHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.GetDB()
		role := c.MustGet("role").(middlewares.AuthRole)

		switch role {
		case middlewares.Admin:
			{
				request := ChangeStatusPhotoRequest{}
				if err := c.BindJSON(&request); err != nil {
					response := gin.H{"Message": "ERROR: BAD REQUEST"}
					c.JSON(http.StatusBadRequest, response)
					return
				}

				query := ChangeStatusPhotoQuery{}
				if err := c.BindQuery(&query); err != nil {
					response := gin.H{"Message": "ERROR: BAD REQUEST"}
					c.JSON(http.StatusBadRequest, response)
					return
				}

				adminID := c.MustGet("id").(uint)
				oldPhoto := models.Photo{Model: gorm.Model{ID: query.PhotoID}}
				newPhoto := models.Photo{Status: request.Status, AdminID: adminID}
				if err := db.Where(&oldPhoto).Updates(&newPhoto).Error; err != nil {
					response := gin.H{"Message": "ERROR: BAD REQUEST"}
					c.JSON(http.StatusBadRequest, response)
					return
				}

				response := gin.H{"Message": "SUCCESS"}
				c.JSON(http.StatusOK, response)
				return
			}
		default:
			{
				response := gin.H{"Message": "ERROR: INVALID ROLE"}
				c.JSON(http.StatusUnauthorized, response)
				return
			}
		}
	}
}

func DeletePhotoHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.GetDB()
		role := c.MustGet("role").(middlewares.AuthRole)

		switch role {
		case middlewares.Team:
			{
				request := DeletePhotoRequest{}
				if err := c.BindJSON(&request); err != nil {
					response := gin.H{"Message": "ERROR: BAD REQUEST"}
					c.JSON(http.StatusBadRequest, response)
					return
				}

				teamID := c.MustGet("id").(uint)
				conditionMembership := models.Membership{TeamID: teamID, ParticipantID: request.ParticipantID}
				membership := models.Membership{}
				if err := db.Where(&conditionMembership).Find(&membership).Error; err != nil {
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
				conditionPhoto := models.Photo{FileName: fileUUID}
				photo := models.Photo{}
				if err := db.Where(&conditionPhoto).Delete(&photo).Error; err != nil {
					response := gin.H{"Message": "ERROR: BAD REQUEST"}
					c.JSON(http.StatusBadRequest, response)
					return
				}

				response := gin.H{"Message": "SUCCESS"}
				c.JSON(http.StatusOK, response)
				return
			}
		default:
			{
				response := gin.H{"Message": "ERROR: INVALID ROLE"}
				c.JSON(http.StatusUnauthorized, response)
				return
			}
		}
	}
}
