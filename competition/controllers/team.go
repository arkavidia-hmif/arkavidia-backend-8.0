package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"

	authenticationConfig "arkavidia-backend-8.0/competition/config/authentication"
	"arkavidia-backend-8.0/competition/middlewares"
	"arkavidia-backend-8.0/competition/models"
	databaseService "arkavidia-backend-8.0/competition/services/database"
)

type SignInRequest struct {
	Username string `json:"username"`
	Password []byte `json:"password"`
}

type SignUpRequest struct {
	Username          string   `json:"username"`
	Password          []byte   `json:"password"`
	ListOfMemberEmail []string `json:"email_list"`
}

func SignInHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.GetDB()
		config := authenticationConfig.GetAuthConfig()
		request := SignInRequest{}

		err := c.BindJSON(&request)
		if err != nil {
			response := gin.H{"Message": "Error: Bad Request!"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		team := models.Team{Username: request.Username}
		err = db.Find(&team).Error
		if err != nil {
			response := gin.H{"Message": "Error: Database Experiencing Problems"}
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		err = bcrypt.CompareHashAndPassword(team.HashedPassword, request.Password)
		if err != nil {
			response := gin.H{"Message": "Error: Invalid Username or Password"}
			c.JSON(http.StatusUnauthorized, response)
			return
		}

		authClaims := middlewares.AuthClaims{
			StandardClaims: jwt.StandardClaims{
				Issuer:    config.ApplicationName,
				ExpiresAt: time.Now().Add(config.LoginExpirationDuration).Unix(),
			},
			TeamID: team.ID,
		}
		unsignedAuthToken := jwt.NewWithClaims(config.JWTSigningMethod, authClaims)
		signedAuthToken, err := unsignedAuthToken.SignedString(config.JWTSignatureKey)
		if err != nil {
			response := gin.H{"Message": "Error: JWT Signing Errors"}
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		authTokenString, err := json.Marshal(gin.H{"token": signedAuthToken})
		if err != nil {
			response := gin.H{"Message": "Error: JWT Signing Errors"}
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		response := gin.H{"Message": "Success", "Data": authTokenString}
		c.JSON(http.StatusOK, response)
		return
	}
}
