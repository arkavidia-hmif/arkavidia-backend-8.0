package controllers

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/lib/pq"
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
	Name           string                `json:"name"`
	Email          string                `json:"email"`
	CareerInterest pq.StringArray        `json:"career_interest"`
	Role           models.MembershipRole `json:"role"`
}

type SignUpRequest struct {
	Username string   `json:"username"`
	Password []byte   `json:"password"`
	TeamName string   `json:"team_name"`
	Members  []Member `json:"member_list"`
}

type ChangePasswordRequest struct {
	Password []byte `json:"password"`
}

type CompetitionRegistrationQuery struct {
	TeamCategory models.TeamCategory `form:"competition" field:"competition"`
}

// SLOW RESPONSE
func SignInHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.GetDB()
		config := authenticationConfig.GetAuthConfig()

		request := SignInRequest{}
		if err := c.BindJSON(&request); err != nil {
			response := gin.H{"Message": "ERROR: BAD REQUEST"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		condition := models.Team{Username: request.Username}
		team := models.Team{}
		if err := db.Where(&condition).Find(&team).Error; err != nil {
			response := gin.H{"Message": "ERROR: INVALID USERNAME OR PASSWORD"}
			c.JSON(http.StatusUnauthorized, response)
			return
		}

		if err := bcrypt.CompareHashAndPassword(team.HashedPassword, request.Password); err != nil {
			response := gin.H{"Message": "ERROR: INVALID USERNAME OR PASSWORD"}
			c.JSON(http.StatusUnauthorized, response)
			return
		}

		authClaims := middlewares.AuthClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    config.ApplicationName,
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.LoginExpirationDuration)),
			},
			TeamID: team.ID,
		}

		unsignedAuthToken := jwt.NewWithClaims(config.JWTSigningMethod, authClaims)
		signedAuthToken, err := unsignedAuthToken.SignedString(config.JWTSignatureKey)
		if err != nil {
			response := gin.H{"Message": "ERROR: JWT SIGNING ERROR"}
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		response := gin.H{"Message": "SUCCESS", "Data": signedAuthToken}
		c.JSON(http.StatusCreated, response)
	}
}

// SLOW RESPONSE
func SignUpHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.GetDB()
		config := authenticationConfig.GetAuthConfig()

		request := SignUpRequest{}
		if err := c.BindJSON(&request); err != nil {
			response := gin.H{"Message": "ERROR: BAD REQUEST"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword(request.Password, rand.Intn(bcrypt.MaxCost-bcrypt.MinCost)+bcrypt.MinCost)
		if err != nil {
			response := gin.H{"Message": "ERROR: BCRYPT ERROR"}
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		team := models.Team{Username: request.Username, HashedPassword: hashedPassword, TeamName: request.TeamName}
		if err := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(&team).Error; err != nil {
				return err
			}
			for _, member := range request.Members {
				condition := models.Participant{Name: member.Name, Email: member.Email, CareerInterest: member.CareerInterest}
				participant := models.Participant{}
				if err := tx.FirstOrCreate(&participant, &condition).Error; err != nil {
					return err
				}
				membership := models.Membership{TeamID: team.ID, ParticipantID: participant.ID, Role: member.Role}
				if err := tx.Create(&membership).Error; err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			response := gin.H{"Message": "ERROR: BAD REQUEST"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		authClaims := middlewares.AuthClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    config.ApplicationName,
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.LoginExpirationDuration)),
			},
			TeamID: team.ID,
		}
		unsignedAuthToken := jwt.NewWithClaims(config.JWTSigningMethod, authClaims)
		signedAuthToken, err := unsignedAuthToken.SignedString(config.JWTSignatureKey)
		if err != nil {
			response := gin.H{"Message": "ERROR: JWT SIGNING ERROR"}
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		authTokenString, err := json.Marshal(gin.H{"token": signedAuthToken})
		if err != nil {
			response := gin.H{"Message": "ERROR: JWT SIGNING ERROR"}
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		response := gin.H{"Message": "SUCCESS", "Data": authTokenString}
		c.JSON(http.StatusCreated, response)
	}
}

func GetTeamHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.GetDB()
		teamID := c.MustGet("team_id").(uint)

		condition := models.Team{Model: gorm.Model{ID: teamID}}
		team := models.Team{}
		if err := db.Preload("Memberships").Where(&condition).Find(&team).Error; err != nil {
			response := gin.H{"Message": "ERROR: BAD REQUEST"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		response := gin.H{"Message": "SUCCESS", "Data": team}
		c.JSON(http.StatusOK, response)
	}
}

// SLOW RESPONSE
func ChangePasswordHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.GetDB()
		teamID := c.MustGet("team_id").(uint)

		request := ChangePasswordRequest{}
		hashedPassword, err := bcrypt.GenerateFromPassword(request.Password, rand.Intn(bcrypt.MaxCost-bcrypt.MinCost)+bcrypt.MinCost)
		if err != nil {
			response := gin.H{"Message": "ERROR: BCRYPT ERROR"}
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		oldTeam := models.Team{Model: gorm.Model{ID: teamID}}
		newTeam := models.Team{HashedPassword: hashedPassword}
		if err := db.Where(&oldTeam).Updates(&newTeam).Error; err != nil {
			response := gin.H{"Message": "ERROR: BAD REQUEST"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		response := gin.H{"Message": "SUCCESS"}
		c.JSON(http.StatusOK, response)
	}
}

func CompetitionRegistration() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.GetDB()
		teamID := c.MustGet("team_id").(uint)

		query := CompetitionRegistrationQuery{}
		if err := c.BindQuery(&query); err != nil {
			response := gin.H{"Message": "ERROR: BAD REQUEST"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		oldTeam := models.Team{Model: gorm.Model{ID: teamID}}
		newTeam := models.Team{TeamCategory: query.TeamCategory}
		if err := db.Where(&oldTeam).Updates(&newTeam).Error; err != nil {
			response := gin.H{"Message": "ERROR: BAD REQUEST"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		response := gin.H{"Message": "SUCCESS"}
		c.JSON(http.StatusOK, response)
	}
}
