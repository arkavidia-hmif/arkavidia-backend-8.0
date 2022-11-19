package controllers

import (
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
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
	TeamName     string   `json:"team_name"`
	ListOfMember []Member `json:"member_list"`
}

type ChangePasswordRequest struct {
	Password []byte `json:"password"`
}

type CompetitionRegistrationQuery struct {
	TeamCategory models.TeamCategory `form:"competition" field:"competition"`
}

func SignInHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.GetDB()
		config := authenticationConfig.GetAuthConfig()
		request := SignInRequest{}

		if err := c.BindJSON(&request); err != nil {
			response := gin.H{"Message": "Error: Bad Request!"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		team := models.Team{Username: request.Username}
		if err := db.Find(&team).Error; err != nil {
			response := gin.H{"Message": "Error: Bad Request!"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		if err := bcrypt.CompareHashAndPassword(team.HashedPassword, request.Password); err != nil {
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

func SignUpHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.GetDB()
		config := authenticationConfig.GetAuthConfig()
		request := SignUpRequest{}

		if err := c.BindJSON(&request); err != nil {
			response := gin.H{"Message": "Error: Bad Request!"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword(request.Password, rand.Intn(bcrypt.MaxCost-bcrypt.MinCost)+bcrypt.MinCost)
		if err != nil {
			response := gin.H{"Message": "Error: Bcrypt Errors"}
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		team := models.Team{Username: request.Username, HashedPassword: hashedPassword, TeamName: request.TeamName}

		if err := db.Transaction(func(tx *gorm.DB) error {
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
		}); err != nil {
			response := gin.H{"Message": "Error: Bad Request!"}
			c.JSON(http.StatusBadRequest, response)
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

func GetTeam() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.GetDB()
		teamID := c.MustGet("team_id").(uuid.UUID)

		team := models.Team{ID: teamID}
		if err := db.Find(&team).Error; err != nil {
			response := gin.H{"Message": "Error: Bad Request!"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		response := gin.H{"Message": "Success", "Data": team}
		c.JSON(http.StatusOK, response)
		return
	}
}

func ChangePasswordHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.GetDB()
		teamID := c.MustGet("team_id").(uuid.UUID)
		request := ChangePasswordRequest{}

		hashedPassword, err := bcrypt.GenerateFromPassword(request.Password, rand.Intn(bcrypt.MaxCost-bcrypt.MinCost)+bcrypt.MinCost)
		if err != nil {
			response := gin.H{"Message": "Error: Bcrypt Errors"}
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		oldTeam := models.Team{ID: teamID}
		newTeam := models.Team{HashedPassword: hashedPassword}
		if err := db.Find(&oldTeam).Updates(&newTeam); err != nil {
			response := gin.H{"Message": "Error: Bad Request!"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		response := gin.H{"Message": "Success"}
		c.JSON(http.StatusOK, response)
		return
	}
}

func CompetitionRegistration() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.GetDB()
		teamID := c.MustGet("team_id").(uuid.UUID)
		query := CompetitionRegistrationQuery{}

		if err := c.BindQuery(&query); err != nil {
			response := gin.H{"Message": "Error: Bad Request!"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		oldTeam := models.Team{ID: teamID}
		newTeam := models.Team{TeamCategory: query.TeamCategory}
		if err := db.Find(&oldTeam).Updates(&newTeam); err != nil {
			response := gin.H{"Message": "Error: Bad Request!"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		response := gin.H{"Message": "Success"}
		c.JSON(http.StatusOK, response)
		return
	}
}
