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

type GetSubmissionQuery struct {
	TeamID uint `form:"team_id" field:"team_id" binding:"required"`
}

type GetAllSubmissionsQuery struct {
	Page int `form:"page" field:"page" binding:"required"`
	Size int `form:"size" field:"size" binding:"required"`
}

type DownloadSubmissionQuery struct {
	SubmissionID uint `form:"submission_id" field:"submission_id" binding:"required"`
}

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
		role := c.MustGet("role").(middlewares.AuthRole)

		switch role {
		case middlewares.Admin:
			{
				query := GetSubmissionQuery{}
				if err := c.BindQuery(&query); err != nil {
					response := gin.H{"Message": "ERROR: BAD REQUEST"}
					c.JSON(http.StatusBadRequest, response)
					return
				}

				condition := models.Submission{TeamID: query.TeamID}
				submissions := []models.Submission{}
				if err := db.Where(&condition).Find(&submissions).Error; err != nil {
					response := gin.H{"Message": "ERROR: BAD REQUEST"}
					c.JSON(http.StatusBadRequest, response)
					return
				}

				response := gin.H{"Message": "SUCCESS", "Data": submissions, "URL": fmt.Sprintf("%s/%s/%s/", config.StorageHost, config.BucketName, config.SubmissionDir)}
				c.JSON(http.StatusOK, response)
				return
			}
		case middlewares.Team:
			{
				teamID := c.MustGet("id").(uint)
				condition := models.Submission{TeamID: teamID}
				submissions := []models.Submission{}
				if err := db.Where(&condition).Find(&submissions).Error; err != nil {
					response := gin.H{"Message": "ERROR: BAD REQUEST"}
					c.JSON(http.StatusBadRequest, response)
					return
				}

				response := gin.H{"Message": "SUCCESS", "Data": submissions, "URL": fmt.Sprintf("%s/%s/%s/", config.StorageHost, config.BucketName, config.SubmissionDir)}
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

func GetAllSubmissionsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.GetDB()
		role := c.MustGet("role").(middlewares.AuthRole)

		switch role {
		case middlewares.Admin:
			{
				query := GetAllSubmissionsQuery{}
				if err := c.BindQuery(&query); err != nil {
					response := gin.H{"Message": "ERROR: BAD REQUEST"}
					c.JSON(http.StatusBadRequest, response)
					return
				}

				offset := (query.Page - 1) * query.Size
				limit := query.Size
				submissions := []models.Submission{}

				if err := db.Offset(offset).Limit(limit).Find(&submissions).Error; err != nil {
					response := gin.H{"Message": "ERROR: BAD REQUEST"}
					c.JSON(http.StatusBadRequest, response)
					return
				}

				response := gin.H{"Message": "SUCCESS", "Data": submissions}
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

func DownloadSubmissionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.GetDB()
		client := storageService.GetClient()
		config := storageConfig.GetStorageConfig()
		role := c.MustGet("role").(middlewares.AuthRole)

		switch role {
		case middlewares.Admin:
			{
				query := DownloadSubmissionQuery{}
				if err := c.BindQuery(&query); err != nil {
					response := gin.H{"Message": "ERROR: BAD REQUEST"}
					c.JSON(http.StatusBadRequest, response)
					return
				}

				conditionSubmission := models.Submission{Model: gorm.Model{ID: query.SubmissionID}}
				submission := models.Submission{}
				if err := db.Where(&conditionSubmission).Find(&submission).Error; err != nil {
					response := gin.H{"Message": "ERROR: BAD REQUEST"}
					c.JSON(http.StatusBadRequest, response)
					return
				}

				filename := fmt.Sprintf("%s.%s", submission.FileName, submission.FileExtension)
				IOWriter, err := storageService.DownloadFile(client, filename, config.SubmissionDir)
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

func AddSubmissionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.GetDB()
		client := storageService.GetClient()
		config := storageConfig.GetStorageConfig()
		role := c.MustGet("role").(middlewares.AuthRole)

		switch role {
		case middlewares.Team:
			{
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
				defer openedFile.Close()

				fileUUID := uuid.New()
				fileExt := filepath.Ext(request.File.Filename)

				teamID := c.MustGet("id").(uint)
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

func DeleteSubmissionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.GetDB()
		role := c.MustGet("role").(middlewares.AuthRole)

		switch role {
		case middlewares.Team:
			{
				teamID := c.MustGet("id").(uint)

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
