package controllers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"gorm.io/gorm"

	storageConfig "arkavidia-backend-8.0/competition/config/storage"
	"arkavidia-backend-8.0/competition/middlewares"
	"arkavidia-backend-8.0/competition/models"
	"arkavidia-backend-8.0/competition/repository"
	databaseService "arkavidia-backend-8.0/competition/services/database"
	storageService "arkavidia-backend-8.0/competition/services/storage"
)

func GetSubmissionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.DB.GetConnection()
		config := storageConfig.Config.GetMetadata()
		response := repository.Response[[]models.Submission]{}

		value, exists := c.Get("role")
		if !exists {
			response.Message = "UNAUTHORIZED"
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		role := value.(middlewares.AuthRole)

		switch role {
		case middlewares.Admin:
			{
				query := repository.GetSubmissionQuery{}
				if err := c.ShouldBindQuery(&query); err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				condition := models.Submission{TeamID: query.TeamID}
				submissions := []models.Submission{}
				if err := db.Where(&condition).Find(&submissions).Error; err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				response.Message = "SUCCESS"
				response.Data = submissions
				response.URL = fmt.Sprintf("%s/%s/%s/", config.StorageHost, config.BucketName, config.SubmissionDir)
				c.JSON(http.StatusOK, response)
				return
			}
		case middlewares.Team:
			{
				value, exists := c.Get("id")
				if !exists {
					response.Message = "UNAUTHORIZED"
					c.AbortWithStatusJSON(http.StatusUnauthorized, response)
					return
				}

				teamID := value.(uint)
				condition := models.Submission{TeamID: teamID}
				submissions := []models.Submission{}
				if err := db.Where(&condition).Find(&submissions).Error; err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				response.Message = "SUCCESS"
				response.Data = submissions
				response.URL = fmt.Sprintf("%s/%s/%s/", config.StorageHost, config.BucketName, config.SubmissionDir)
				c.JSON(http.StatusOK, response)
				return
			}
		default:
			{
				response.Message = "ERROR: INVALID ROLE"
				c.AbortWithStatusJSON(http.StatusUnauthorized, response)
				return
			}
		}
	}
}

func GetAllSubmissionsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.DB.GetConnection()
		config := storageConfig.Config.GetMetadata()
		response := repository.Response[[]models.Submission]{}

		value, exists := c.Get("role")
		if !exists {
			response.Message = "UNAUTHORIZED"
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		role := value.(middlewares.AuthRole)

		switch role {
		case middlewares.Admin:
			{
				query := repository.GetAllSubmissionsQuery{}
				if err := c.ShouldBindQuery(&query); err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				offset := (query.Page - 1) * query.Size
				limit := query.Size
				submissions := []models.Submission{}

				if err := db.Offset(offset).Limit(limit).Find(&submissions).Error; err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				response.Message = "SUCCESS"
				response.Data = submissions
				response.URL = fmt.Sprintf("%s/%s/%s/", config.StorageHost, config.BucketName, config.SubmissionDir)
				c.JSON(http.StatusOK, response)
				return
			}
		default:
			{
				response.Message = "ERROR: INVALID ROLE"
				c.AbortWithStatusJSON(http.StatusUnauthorized, response)
				return
			}
		}
	}
}

func DownloadSubmissionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.DB.GetConnection()
		config := storageConfig.Config.GetMetadata()
		response := repository.Response[models.Submission]{}

		value, exists := c.Get("role")
		if !exists {
			response.Message = "UNAUTHORIZED"
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		role := value.(middlewares.AuthRole)

		switch role {
		case middlewares.Admin:
			{
				query := repository.DownloadSubmissionQuery{}
				if err := c.ShouldBindQuery(&query); err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				condition := models.Submission{Model: gorm.Model{ID: query.SubmissionID}}
				submission := models.Submission{}
				if err := db.Where(&condition).Find(&submission).Error; err != nil {
					response.Message = "ERROR: CONTENT NOT FOUND IN DB"
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				filename := fmt.Sprintf("%s.%s", submission.FileName, submission.FileExtension)
				IOWriter, err := storageService.Client.DownloadFile(filename, config.SubmissionDir)
				if err != nil {
					response.Message = "ERROR: CONTENT NOT FOUND IN STORAGE"
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				var content []byte
				length, err := IOWriter.Write(content)
				if err != nil {
					response.Message = "ERROR: CONTENT CANNOT BE WRITTEN"
					c.AbortWithStatusJSON(http.StatusInternalServerError, response)
					return
				}

				c.Header("Content-Description", "File Transfer")
				c.Header("Content-Transfer-Encoding", "binary")
				c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
				c.Header("Content-Type", "application/octet-stream")
				c.Header("Accept-Length", fmt.Sprintf("%d", length))
				c.Writer.Write(content)

				response.Message = "SUCCESS"
				c.JSON(http.StatusOK, response)
				return
			}
		case middlewares.Team:
			{
				query := repository.DownloadSubmissionQuery{}
				if err := c.ShouldBindQuery(&query); err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				value, exists := c.Get("id")
				if !exists {
					response.Message = "UNAUTHORIZED"
					c.AbortWithStatusJSON(http.StatusUnauthorized, response)
					return
				}

				teamID := value.(uint)
				condition := models.Submission{Model: gorm.Model{ID: query.SubmissionID}, TeamID: teamID}
				submission := models.Submission{}
				if err := db.Where(&condition).Find(&submission).Error; err != nil {
					response.Message = "ERROR: CONTENT NOT FOUND IN DB"
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				filename := fmt.Sprintf("%s.%s", submission.FileName, submission.FileExtension)
				IOWriter, err := storageService.Client.DownloadFile(filename, config.SubmissionDir)
				if err != nil {
					response.Message = "ERROR: CONTENT NOT FOUND IN STORAGE"
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				var content []byte
				length, err := IOWriter.Write(content)
				if err != nil {
					response.Message = "ERROR: CONTENT CANNOT BE WRITTEN"
					c.AbortWithStatusJSON(http.StatusInternalServerError, response)
					return
				}

				c.Header("Content-Description", "File Transfer")
				c.Header("Content-Transfer-Encoding", "binary")
				c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
				c.Header("Content-Type", "application/octet-stream")
				c.Header("Accept-Length", fmt.Sprintf("%d", length))
				c.Writer.Write(content)

				response.Message = "SUCCESS"
				c.JSON(http.StatusOK, response)
				return
			}
		default:
			{
				response.Message = "ERROR: INVALID ROLE"
				c.AbortWithStatusJSON(http.StatusUnauthorized, response)
				return
			}
		}
	}
}

func RenderSubmissionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.DB.GetConnection()
		config := storageConfig.Config.GetMetadata()
		response := repository.Response[models.Submission]{}

		value, exists := c.Get("role")
		if !exists {
			response.Message = "UNAUTHORIZED"
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		role := value.(middlewares.AuthRole)

		switch role {
		case middlewares.Admin:
			{
				query := repository.DownloadSubmissionQuery{}
				if err := c.ShouldBindQuery(&query); err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				condition := models.Submission{Model: gorm.Model{ID: query.SubmissionID}}
				submission := models.Submission{}
				if err := db.Where(&condition).Find(&submission).Error; err != nil {
					response.Message = "ERROR: CONTENT NOT FOUND IN DB"
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				url := "https://storage.googleapis.com/arkavidia-8/competition/submission/23f672bb-7b67-4759-bc6b-17783494b208.pdf"
				// url := fmt.Sprintf("%s/%s/%s/%s%s", config.StorageHost, config.BucketName, config.SubmissionDir, submission.FileName, submission.FileExtension)
				res, err := http.Get(url)
				if err != nil {
					response.Message = err.Error()
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}
				defer res.Body.Close()

				content, err := ioutil.ReadAll(res.Body)
				if err != nil {
					response.Message = "ERROR: CONTENT CANNOT BE WRITTEN"
					c.AbortWithStatusJSON(http.StatusInternalServerError, response)
					return
				}

				mtype, err := mimetype.DetectReader(bytes.NewReader(content))
				if err != nil {
					response.Message = "ERROR: CANNOT GET CONTENT TYPE"
				}

				c.Header("Content-Description", "File Transfer")
				c.Header("Content-Transfer-Encoding", "binary")
				c.Header("Content-Disposition", "inline")
				c.Header("Content-Type", mtype.String())
				c.Header("Accept-Length", fmt.Sprintf("%d", res.ContentLength))
				c.Writer.Write(content)

				response.Message = "SUCCESS"
				c.JSON(http.StatusOK, response)
				return
			}
		case middlewares.Team:
			{
				query := repository.DownloadSubmissionQuery{}
				if err := c.ShouldBindQuery(&query); err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				value, exists := c.Get("id")
				if !exists {
					response.Message = "UNAUTHORIZED"
					c.AbortWithStatusJSON(http.StatusUnauthorized, response)
					return
				}

				teamID := value.(uint)
				condition := models.Submission{Model: gorm.Model{ID: query.SubmissionID}, TeamID: teamID}
				submission := models.Submission{}
				if err := db.Where(&condition).Find(&submission).Error; err != nil {
					response.Message = "ERROR: CONTENT NOT FOUND IN DB"
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				url := fmt.Sprintf("%s/%s/%s/%s%s", config.StorageHost, config.BucketName, config.SubmissionDir, submission.FileName, submission.FileExtension)
				res, err := http.Get(url)
				if err != nil {
					response.Message = err.Error()
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}
				defer res.Body.Close()

				content, err := ioutil.ReadAll(res.Body)
				if err != nil {
					response.Message = "ERROR: CONTENT CANNOT BE WRITTEN"
					c.AbortWithStatusJSON(http.StatusInternalServerError, response)
					return
				}

				mtype, err := mimetype.DetectReader(bytes.NewReader(content))
				if err != nil {
					response.Message = "ERROR: CANNOT GET CONTENT TYPE"
				}

				c.Header("Content-Description", "File Transfer")
				c.Header("Content-Transfer-Encoding", "binary")
				c.Header("Content-Disposition", "inline")
				c.Header("Content-Type", mtype.String())
				c.Header("Accept-Length", fmt.Sprintf("%d", res.ContentLength))
				c.Writer.Write(content)

				response.Message = "SUCCESS"
				c.JSON(http.StatusOK, response)
				return
			}
		default:
			{
				response.Message = "ERROR: INVALID ROLE"
				c.AbortWithStatusJSON(http.StatusUnauthorized, response)
				return
			}
		}
	}
}

func AddSubmissionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.DB.GetConnection()
		config := storageConfig.Config.GetMetadata()
		response := repository.Response[models.Submission]{}

		value, exists := c.Get("role")
		if !exists {
			response.Message = "UNAUTHORIZED"
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		role := value.(middlewares.AuthRole)

		switch role {
		case middlewares.Team:
			{
				request := repository.AddSubmissionRequest{}
				if err := c.ShouldBindWith(&request, binding.FormMultipart); err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				openedFile, err := request.File.Open()
				if err != nil {
					response.Message = "ERROR: FILE CANNOT BE ACCESSED"
					c.AbortWithStatusJSON(http.StatusInternalServerError, response)
					return
				}
				defer openedFile.Close()

				fileUUID := uuid.New()
				fileExt := filepath.Ext(request.File.Filename)

				value, exists := c.Get("id")
				if !exists {
					response.Message = "UNAUTHORIZED"
					c.AbortWithStatusJSON(http.StatusUnauthorized, response)
					return
				}

				teamID := value.(uint)
				submission := models.Submission{FileName: fileUUID, FileExtension: fileExt, TeamID: teamID, Stage: request.Stage}
				if err := db.Create(&submission).Error; err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				if err := storageService.Client.UploadFile(fmt.Sprintf("%s%s", fileUUID, fileExt), config.SubmissionDir, openedFile); err != nil {
					response.Message = "ERROR: GOOGLE CLOUD STORAGE CANNOT BE ACCESSED"
					c.AbortWithStatusJSON(http.StatusInternalServerError, response)
					return
				}

				response.Message = "SUCCESS"
				response.Data = submission
				response.URL = fmt.Sprintf("%s/%s/%s/", config.StorageHost, config.BucketName, config.SubmissionDir)
				c.JSON(http.StatusCreated, response)
				return
			}
		default:
			{
				response.Message = "ERROR: INVALID ROLE"
				c.AbortWithStatusJSON(http.StatusUnauthorized, response)
				return
			}
		}
	}
}

func DeleteSubmissionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.DB.GetConnection()
		response := repository.Response[models.Submission]{}

		value, exists := c.Get("role")
		if !exists {
			response.Message = "UNAUTHORIZED"
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		role := value.(middlewares.AuthRole)

		switch role {
		case middlewares.Team:
			{
				request := repository.DeleteSubmissionRequest{}
				if err := c.ShouldBindJSON(&request); err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				fileUUID, err := uuid.Parse(request.FileName)
				if err != nil {
					response.Message = "ERROR: INVALID FILENAME"
					c.AbortWithStatusJSON(http.StatusInternalServerError, response)
					return
				}

				value, exists := c.Get("id")
				if !exists {
					response.Message = "UNAUTHORIZED"
					c.AbortWithStatusJSON(http.StatusUnauthorized, response)
					return
				}

				teamID := value.(uint)
				condition := models.Submission{FileName: fileUUID, TeamID: teamID}
				submission := models.Submission{}
				if err := db.Where(&condition).Delete(&submission).Error; err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				response.Message = "SUCCESS"
				c.JSON(http.StatusOK, response)
				return
			}
		default:
			{
				response.Message = "ERROR: INVALID ROLE"
				c.AbortWithStatusJSON(http.StatusUnauthorized, response)
				return
			}
		}
	}
}
