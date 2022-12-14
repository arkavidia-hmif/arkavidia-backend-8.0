package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"

	authConfig "arkavidia-backend-8.0/competition/config/authentication"
	"arkavidia-backend-8.0/competition/middlewares"
	"arkavidia-backend-8.0/competition/models"
	"arkavidia-backend-8.0/competition/repository"
	databaseService "arkavidia-backend-8.0/competition/services/database"
)

func SignInAdminHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.DB.GetConnection()
		config := authConfig.Config.GetMetadata()
		response := repository.Response[string]{}

		request := repository.SignInAdminRequest{}
		if err := c.ShouldBindJSON(&request); err != nil {
			response.Message = "ERROR: BAD REQUEST"
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		condition := models.Admin{Username: request.Username}
		admin := models.Admin{}
		if err := db.Where(&condition).Find(&admin).Error; err != nil {
			response.Message = "ERROR: INVALID USERNAME"
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		if err := bcrypt.CompareHashAndPassword(admin.HashedPassword, []byte(request.Password)); err != nil {
			response.Message = "ERROR: INVALID PASSWORD"
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		adminClaims := middlewares.AuthClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    config.ApplicationName,
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.LoginExpirationDuration)),
			},
			ID:   admin.ID,
			Role: middlewares.Admin,
		}

		unsignedAuthToken := jwt.NewWithClaims(config.JWTSigningMethod, adminClaims)
		signedAuthToken, err := unsignedAuthToken.SignedString(config.JWTSignatureKey)
		if err != nil {
			response.Message = "ERROR: JWT SIGNING ERROR"
			c.AbortWithStatusJSON(http.StatusInternalServerError, response)
			return
		}

		response.Message = "SUCCESS"
		response.Data = signedAuthToken
		c.JSON(http.StatusCreated, response)
	}
}
