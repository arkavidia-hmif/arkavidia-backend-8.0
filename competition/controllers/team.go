package controllers

import (
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	authenticationConfig "arkavidia-backend-8.0/competition/config/authentication"
	"arkavidia-backend-8.0/competition/middlewares"
	"arkavidia-backend-8.0/competition/models"
	databaseService "arkavidia-backend-8.0/competition/services/database"
)

type SignInRequest struct {
	Username string `json:"username"`
	Password []byte `json:"password"`
}

type Member struct {
	Email string
	Name  string
	Role  models.MembershipRole
}

type SignUpRequest struct {
	Username     string   `json:"username"`
	Password     []byte   `json:"password"`
	ListOfMember []Member `json:"member_list"`
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

		response := gin.H{"Message": "Success", "Data": authTokenString} // TODO: Return value baru token saja
		c.JSON(http.StatusOK, response)
		return
	}
}

func SignUpHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.GetDB()
		config := authenticationConfig.GetAuthConfig()
		request := SignUpRequest{}

		err := c.BindJSON(&request)
		if err != nil {
			response := gin.H{"Message": "Error: Bad Request!"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword(request.Password, rand.Intn(bcrypt.MaxCost-bcrypt.MinCost)+bcrypt.MinCost)
		team := models.Team{Username: request.Username, HashedPassword: hashedPassword}

		err = db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(&team).Error; err != nil {
				return err
			}
			for _, member := range request.ListOfMember {
				participant := models.Participant{Name: member.Name, Email: member.Email}
				if err := tx.Find(&participant).Error; err != nil {
					if !errors.Is(err, gorm.ErrRecordNotFound) {
						return err
					}
				}
				if err := tx.Create(&participant).Error; err != nil {
					return err
				}
				membership := models.Membership{TeamID: team.ID, ParticipantID: participant.ID, Role: member.Role}
				if err := tx.Create(&membership).Error; err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			response := gin.H{"Message": "Error: Database Experiencing Problems"}
			c.JSON(http.StatusInternalServerError, response)
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

		response := gin.H{"Message": "Success", "Data": authTokenString} // TODO: Return value baru token saja
		c.JSON(http.StatusOK, response)
		return
	}
}
