package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"

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
	"arkavidia-backend-8.0/competition/types"
	"arkavidia-backend-8.0/competition/utils/sanitizer"
)

func GetPhotoHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.DB.GetConnection()
		config := storageConfig.Config.GetMetadata()
		response := sanitizer.Response[[]models.Photo]{}

		value, exists := c.Get("role")
		if !exists {
			response.Message = "UNAUTHORIZED"
			c.AbortWithStatusJSON(http.StatusUnauthorized, sanitizer.SanitizeStruct(response))
			return
		}

		role := value.(middlewares.AuthRole)

		switch role {
		case middlewares.Admin:
			{
				query := repository.GetPhotoQuery{}
				if err := c.ShouldBindQuery(&query); err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, sanitizer.SanitizeStruct(response))
					return
				}

				conditionPhoto := models.Photo{ParticipantID: query.ParticipantID}
				photos := []models.Photo{}
				if err := db.Where(&conditionPhoto).Find(&photos).Error; err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, sanitizer.SanitizeStruct(response))
					return
				}

				response.Message = "SUCCESS"
				response.Data = photos
				response.URL = fmt.Sprintf("%s/%s/%s", config.StorageHost, config.BucketName, config.PhotoDir)
				c.JSON(http.StatusOK, sanitizer.SanitizeStruct(response))
				return
			}
		case middlewares.Team:
			{
				query := repository.GetPhotoQuery{}
				if err := c.ShouldBindQuery(&query); err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, sanitizer.SanitizeStruct(response))
					return
				}

				value, exists := c.Get("id")
				if !exists {
					response.Message = "UNAUTHORIZED"
					c.AbortWithStatusJSON(http.StatusUnauthorized, sanitizer.SanitizeStruct(response))
					return
				}

				teamID := value.(uint)
				conditionMembership := models.Membership{TeamID: teamID, ParticipantID: query.ParticipantID}
				membership := models.Membership{}
				if err := db.Where(&conditionMembership).Find(&membership).Error; err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, sanitizer.SanitizeStruct(response))
					return
				}

				conditionPhoto := models.Photo{ParticipantID: query.ParticipantID}
				photos := []models.Photo{}
				if err := db.Where(&conditionPhoto).Find(&photos).Error; err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, sanitizer.SanitizeStruct(response))
					return
				}

				response.Message = "SUCCESS"
				response.Data = photos
				response.URL = fmt.Sprintf("%s/%s/%s", config.StorageHost, config.BucketName, config.PhotoDir)
				c.JSON(http.StatusOK, sanitizer.SanitizeStruct(response))
				return
			}
		default:
			{
				response.Message = "ERROR: INVALID ROLE"
				c.AbortWithStatusJSON(http.StatusUnauthorized, sanitizer.SanitizeStruct(response))
				return
			}
		}
	}
}

func GetAllPhotosHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.DB.GetConnection()
		response := sanitizer.Response[[]models.Photo]{}

		value, exists := c.Get("role")
		if !exists {
			response.Message = "UNAUTHORIZED"
			c.AbortWithStatusJSON(http.StatusUnauthorized, sanitizer.SanitizeStruct(response))
			return
		}

		role := value.(middlewares.AuthRole)

		switch role {
		case middlewares.Admin:
			{
				query := repository.GetAllPhotosQuery{}
				if err := c.ShouldBindQuery(&query); err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, sanitizer.SanitizeStruct(response))
					return
				}

				offset := (query.Page - 1) * query.Size
				limit := query.Size
				photos := []models.Photo{}

				if err := db.Offset(offset).Limit(limit).Find(&photos).Error; err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, sanitizer.SanitizeStruct(response))
					return
				}

				response.Message = "SUCCESS"
				response.Data = photos
				c.JSON(http.StatusOK, sanitizer.SanitizeStruct(response))
				return
			}
		default:
			{
				response.Message = "ERROR: INVALID ROLE"
				c.AbortWithStatusJSON(http.StatusUnauthorized, sanitizer.SanitizeStruct(response))
				return
			}
		}
	}
}

func DownloadPhotoHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.DB.GetConnection()
		config := storageConfig.Config.GetMetadata()
		response := sanitizer.Response[models.Photo]{}

		value, exists := c.Get("role")
		if !exists {
			response.Message = "UNAUTHORIZED"
			c.AbortWithStatusJSON(http.StatusUnauthorized, sanitizer.SanitizeStruct(response))
			return
		}

		role := value.(middlewares.AuthRole)

		switch role {
		case middlewares.Admin:
			{
				query := repository.AdminDownloadPhotoQuery{}
				if err := c.ShouldBindQuery(&query); err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, sanitizer.SanitizeStruct(response))
					return
				}

				condition := models.Photo{Model: gorm.Model{ID: query.PhotoID}}
				photo := models.Photo{}
				if err := db.Where(&condition).Find(&photo).Error; err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, sanitizer.SanitizeStruct(response))
					return
				}

				filename := fmt.Sprintf("%s.%s", photo.FileName, photo.FileExtension)
				IOWriter, err := storageService.Client.DownloadFile(filename, config.PhotoDir)
				if err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, sanitizer.SanitizeStruct(response))
					return
				}

				var content []byte
				length, err := IOWriter.Write(content)
				if err != nil {
					response.Message = "ERROR: CONTENT CANNOT BE WRITTEN"
					c.AbortWithStatusJSON(http.StatusInternalServerError, sanitizer.SanitizeStruct(response))
					return
				}

				c.Header("Content-Description", "File Transfer")
				c.Header("Content-Transfer-Encoding", "binary")
				c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
				c.Header("Content-Type", "application/octet-stream")
				c.Header("Accept-Length", fmt.Sprintf("%d", length))
				c.Writer.Write(content)

				response.Message = "SUCCESS"
				c.JSON(http.StatusOK, sanitizer.SanitizeStruct(response))
				return
			}
		case middlewares.Team:
			{
				query := repository.TeamDownloadPhotoQuery{}
				if err := c.ShouldBindQuery(&query); err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, sanitizer.SanitizeStruct(response))
					return
				}

				value, exists := c.Get("id")
				if !exists {
					response.Message = "UNAUTHORIZED"
					c.AbortWithStatusJSON(http.StatusUnauthorized, sanitizer.SanitizeStruct(response))
					return
				}

				teamID := value.(uint)
				conditionMembership := models.Membership{TeamID: teamID, ParticipantID: query.ParticipantID}
				membership := models.Membership{}
				if err := db.Where(&conditionMembership).Find(&membership).Error; err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, sanitizer.SanitizeStruct(response))
					return
				}

				conditionPhoto := models.Photo{Model: gorm.Model{ID: query.PhotoID}}
				photo := models.Photo{}
				if err := db.Where(&conditionPhoto).Find(&photo).Error; err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, sanitizer.SanitizeStruct(response))
					return
				}

				filename := fmt.Sprintf("%s.%s", photo.FileName, photo.FileExtension)
				IOWriter, err := storageService.Client.DownloadFile(filename, config.PhotoDir)
				if err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, sanitizer.SanitizeStruct(response))
					return
				}

				var content []byte
				length, err := IOWriter.Write(content)
				if err != nil {
					response.Message = "ERROR: CONTENT CANNOT BE WRITTEN"
					c.AbortWithStatusJSON(http.StatusInternalServerError, sanitizer.SanitizeStruct(response))
					return
				}

				c.Header("Content-Description", "File Transfer")
				c.Header("Content-Transfer-Encoding", "binary")
				c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
				c.Header("Content-Type", "application/octet-stream")
				c.Header("Accept-Length", fmt.Sprintf("%d", length))
				c.Writer.Write(content)

				response.Message = "SUCCESS"
				c.JSON(http.StatusOK, sanitizer.SanitizeStruct(response))
				return
			}
		default:
			{
				response.Message = "ERROR: INVALID ROLE"
				c.AbortWithStatusJSON(http.StatusUnauthorized, sanitizer.SanitizeStruct(response))
				return
			}
		}
	}
}

func AddPhotoHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.DB.GetConnection()
		config := storageConfig.Config.GetMetadata()
		response := sanitizer.Response[models.Photo]{}

		value, exists := c.Get("role")
		if !exists {
			response.Message = "UNAUTHORIZED"
			c.AbortWithStatusJSON(http.StatusUnauthorized, sanitizer.SanitizeStruct(response))
			return
		}

		role := value.(middlewares.AuthRole)

		switch role {
		case middlewares.Team:
			{
				request := repository.AddPhotoRequest{}
				if err := c.ShouldBindWith(&request, binding.FormMultipart); err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, sanitizer.SanitizeStruct(response))
					return
				}

				value, exists := c.Get("id")
				if !exists {
					response.Message = "UNAUTHORIZED"
					c.AbortWithStatusJSON(http.StatusUnauthorized, sanitizer.SanitizeStruct(response))
					return
				}

				teamID := value.(uint)
				condition := models.Membership{TeamID: teamID, ParticipantID: request.ParticipantID}
				membership := models.Membership{}
				if err := db.Where(&condition).Find(&membership).Error; err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, sanitizer.SanitizeStruct(response))
					return
				}

				openedFile, err := request.File.Open()
				if err != nil {
					response.Message = "ERROR: FILE CANNOT BE ACCESSED"
					c.AbortWithStatusJSON(http.StatusInternalServerError, sanitizer.SanitizeStruct(response))
					return
				}
				defer openedFile.Close()

				fileUUID := uuid.New()
				fileExt := filepath.Ext(request.File.Filename)

				photo := models.Photo{FileName: fileUUID, FileExtension: fileExt, ParticipantID: request.ParticipantID, Status: types.WaitingForApproval, Type: request.Type}
				if err := db.Create(&photo).Error; err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, sanitizer.SanitizeStruct(response))
					return
				}

				if err := storageService.Client.UploadFile(fmt.Sprintf("%s%s", fileUUID, fileExt), config.PhotoDir, openedFile); err != nil {
					response.Message = "ERROR: GOOGLE CLOUD STORAGE CANNOT BE ACCESSED"
					c.AbortWithStatusJSON(http.StatusInternalServerError, sanitizer.SanitizeStruct(response))
					return
				}

				response.Message = "SUCCESS"
				response.Data = photo
				response.URL = fmt.Sprintf("%s/%s/%s", config.StorageHost, config.BucketName, config.PhotoDir)
				c.JSON(http.StatusCreated, sanitizer.SanitizeStruct(response))
				return
			}
		default:
			{
				response.Message = "ERROR: INVALID ROLE"
				c.AbortWithStatusJSON(http.StatusUnauthorized, sanitizer.SanitizeStruct(response))
				return
			}
		}
	}
}

func ChangeStatusPhotoHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.DB.GetConnection()
		response := sanitizer.Response[models.Photo]{}

		value, exists := c.Get("role")
		if !exists {
			response.Message = "UNAUTHORIZED"
			c.AbortWithStatusJSON(http.StatusUnauthorized, sanitizer.SanitizeStruct(response))
			return
		}

		role := value.(middlewares.AuthRole)

		switch role {
		case middlewares.Admin:
			{
				request := repository.ChangeStatusPhotoRequest{}
				if err := c.ShouldBindJSON(&request); err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, sanitizer.SanitizeStruct(response))
					return
				}

				query := repository.ChangeStatusPhotoQuery{}
				if err := c.ShouldBindQuery(&query); err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, sanitizer.SanitizeStruct(response))
					return
				}

				value, exists := c.Get("id")
				if !exists {
					response.Message = "UNAUTHORIZED"
					c.AbortWithStatusJSON(http.StatusUnauthorized, sanitizer.SanitizeStruct(response))
					return
				}

				adminID := value.(uint)
				oldPhoto := models.Photo{Model: gorm.Model{ID: query.PhotoID}}
				newPhoto := models.Photo{Status: request.Status, AdminID: adminID}
				if err := db.Where(&oldPhoto).Updates(&newPhoto).Error; err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, sanitizer.SanitizeStruct(response))
					return
				}

				response.Message = "SUCCESS"
				c.JSON(http.StatusOK, sanitizer.SanitizeStruct(response))
				return
			}
		default:
			{
				response.Message = "ERROR: INVALID ROLE"
				c.AbortWithStatusJSON(http.StatusUnauthorized, sanitizer.SanitizeStruct(response))
				return
			}
		}
	}
}

func DeletePhotoHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.DB.GetConnection()
		response := sanitizer.Response[models.Photo]{}

		value, exists := c.Get("role")
		if !exists {
			response.Message = "UNAUTHORIZED"
			c.AbortWithStatusJSON(http.StatusUnauthorized, sanitizer.SanitizeStruct(response))
			return
		}

		role := value.(middlewares.AuthRole)

		switch role {
		case middlewares.Team:
			{
				request := repository.DeletePhotoRequest{}
				if err := c.ShouldBindJSON(&request); err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, sanitizer.SanitizeStruct(response))
					return
				}

				value, exists := c.Get("id")
				if !exists {
					response.Message = "UNAUTHORIZED"
					c.AbortWithStatusJSON(http.StatusUnauthorized, sanitizer.SanitizeStruct(response))
					return
				}

				teamID := value.(uint)
				conditionMembership := models.Membership{TeamID: teamID, ParticipantID: request.ParticipantID}
				membership := models.Membership{}
				if err := db.Where(&conditionMembership).Find(&membership).Error; err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, sanitizer.SanitizeStruct(response))
					return
				}

				fileUUID, err := uuid.Parse(request.FileName)
				if err != nil {
					response.Message = "ERROR: INVALID FILENAME"
					c.AbortWithStatusJSON(http.StatusInternalServerError, sanitizer.SanitizeStruct(response))
					return
				}
				conditionPhoto := models.Photo{FileName: fileUUID}
				photo := models.Photo{}
				if err := db.Where(&conditionPhoto).Delete(&photo).Error; err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, sanitizer.SanitizeStruct(response))
					return
				}

				response.Message = "SUCCESS"
				c.JSON(http.StatusOK, sanitizer.SanitizeStruct(response))
				return
			}
		default:
			{
				response.Message = "ERROR: INVALID ROLE"
				c.AbortWithStatusJSON(http.StatusUnauthorized, sanitizer.SanitizeStruct(response))
				return
			}
		}
	}
}
