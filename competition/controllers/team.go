package controllers

import (
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
	"arkavidia-backend-8.0/competition/utils/worker"
)

type SignInTeamRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type SignUpTeamRequest struct {
	Username string             `json:"username" binding:"required"`
	Password string             `json:"password" binding:"required"`
	TeamName string             `json:"team_name" binding:"required"`
	Members  []SignUpMembership `json:"member_list" binding:"required"`
}

type GetTeamQuery struct {
	TeamID uint `form:"team_id" field:"team_id" binding:"required"`
}

type GetAllTeamsQuery struct {
	Page int `form:"page" field:"page" binding:"required"`
	Size int `form:"size" field:"size" binding:"required"`
}

type ChangePasswordRequest struct {
	Password string `json:"password" binding:"required"`
}

type CompetitionRegistrationQuery struct {
	TeamCategory models.TeamCategory `form:"competition" field:"competition" binding:"required"`
}

type ChangeStatusTeamQuery struct {
	TeamID uint `form:"team_id" field:"team_id" binding:"required"`
}

type ChangeStatusTeamRequest struct {
	Status models.TeamStatus `json:"status" binding:"required"`
}

func SignInTeamHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.GetDB()
		config := authenticationConfig.GetAuthConfig()

		request := SignInTeamRequest{}
		if err := c.BindJSON(&request); err != nil {
			response := gin.H{"Message": "ERROR: INCOMPLETE REQUEST"}
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

		if err := bcrypt.CompareHashAndPassword(team.HashedPassword, []byte(request.Password)); err != nil {
			response := gin.H{"Message": "ERROR: INVALID USERNAME OR PASSWORD"}
			c.JSON(http.StatusUnauthorized, response)
			return
		}

		authClaims := middlewares.AuthClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    config.ApplicationName,
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.LoginExpirationDuration)),
			},
			ID:   team.ID,
			Role: middlewares.Team,
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

func SignUpTeamHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.GetDB()
		config := authenticationConfig.GetAuthConfig()

		request := SignUpTeamRequest{}
		if err := c.BindJSON(&request); err != nil {
			response := gin.H{"Message": "ERROR: INCOMPLETE REQUEST"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			response := gin.H{"Message": "ERROR: BCRYPT ERROR"}
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		// Validate Username
		conditionUsername := models.Team{Username: request.Username}
		team := models.Team{}
		if err := db.Where(&conditionUsername).Find(&team).Error; err != nil {
			response := gin.H{"Message": "ERROR: BAD REQUEST"}
			c.JSON(http.StatusBadRequest, response)
			return
		}
		if team.Username != "" {
			response := gin.H{"Message": "ERROR: USERNAME EXISTED"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		// Validate Team Name
		conditionTeamName := models.Team{TeamName: request.TeamName}
		team = models.Team{}
		if err := db.Where(&conditionTeamName).Find(&team).Error; err != nil {
			response := gin.H{"Message": "ERROR: BAD REQUEST"}
			c.JSON(http.StatusBadRequest, response)
			return
		}
		if team.TeamName != "" {
			response := gin.H{"Message": "ERROR: TEAMNAME EXISTED"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		team = models.Team{Username: request.Username, HashedPassword: hashedPassword, TeamName: request.TeamName, Status: models.WaitingForEvaluation}
		if err := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(&team).Error; err != nil {
				return err
			}

			for _, member := range request.Members {
				conditionParticipant := models.Participant{Name: member.Name, Email: member.Email, CareerInterest: member.CareerInterest, Status: models.WaitingForVerification}
				participant := models.Participant{}
				if err := tx.FirstOrCreate(&participant, &conditionParticipant).Error; err != nil {
					return err
				}

				membership := models.Membership{TeamID: team.ID, ParticipantID: participant.ID, Role: member.Role, Team: team, Participant: participant}
				if err := tx.Create(&membership).Error; err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			response := gin.H{"Message": err.Error()}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		authClaims := middlewares.AuthClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    config.ApplicationName,
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.LoginExpirationDuration)),
			},
			ID:   team.ID,
			Role: middlewares.Team,
		}

		unsignedAuthToken := jwt.NewWithClaims(config.JWTSigningMethod, authClaims)
		signedAuthToken, err := unsignedAuthToken.SignedString(config.JWTSignatureKey)
		if err != nil {
			response := gin.H{"Message": "ERROR: JWT SIGNING ERROR"}
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		// Asynchronously mail to every member registered on the team
		for _, member := range request.Members {
			worker.AddMailToBroker(worker.MailParameters{Email: member.Email})
		}

		response := gin.H{"Message": "SUCCESS", "Data": signedAuthToken}
		c.JSON(http.StatusCreated, response)
	}
}

func GetTeamHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.GetDB()
		role := c.MustGet("role").(middlewares.AuthRole)

		switch role {
		case middlewares.Admin:
			{
				query := GetTeamQuery{}
				if err := c.BindQuery(&query); err != nil {
					response := gin.H{"Message": "ERROR: BAD REQUEST"}
					c.JSON(http.StatusBadRequest, response)
					return
				}

				condition := models.Team{Model: gorm.Model{ID: query.TeamID}}
				team := models.Team{}
				if err := db.Preload("Memberships").Where(&condition).Find(&team).Error; err != nil {
					response := gin.H{"Message": "ERROR: BAD REQUEST"}
					c.JSON(http.StatusBadRequest, response)
					return
				}

				response := gin.H{"Message": "SUCCESS", "Data": team}
				c.JSON(http.StatusOK, response)
				return
			}
		case middlewares.Team:
			{
				teamID := c.MustGet("id").(uint)
				condition := models.Team{Model: gorm.Model{ID: teamID}}
				team := models.Team{}
				if err := db.Preload("Memberships").Where(&condition).Find(&team).Error; err != nil {
					response := gin.H{"Message": "ERROR: BAD REQUEST"}
					c.JSON(http.StatusBadRequest, response)
					return
				}

				response := gin.H{"Message": "SUCCESS", "Data": team}
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

func GetAllTeamsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.GetDB()
		role := c.MustGet("role").(middlewares.AuthRole)

		switch role {
		case middlewares.Admin:
			{
				query := GetAllTeamsQuery{}
				if err := c.BindQuery(&query); err != nil {
					response := gin.H{"Message": "ERROR: BAD REQUEST"}
					c.JSON(http.StatusBadRequest, response)
					return
				}

				offset := (query.Page - 1) * query.Size
				limit := query.Size
				teams := []models.Team{}
				if err := db.Offset(offset).Limit(limit).Find(&teams).Error; err != nil {
					response := gin.H{"Message": "ERROR: BAD REQUEST"}
					c.JSON(http.StatusBadRequest, response)
					return
				}

				response := gin.H{"Message": "SUCCESS", "Data": teams}
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

func ChangePasswordHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.GetDB()
		role := c.MustGet("role").(middlewares.AuthRole)

		switch role {
		case middlewares.Team:
			{
				request := ChangePasswordRequest{}
				if err := c.BindJSON(&request); err != nil {
					response := gin.H{"Message": "ERROR: BAD REQUEST"}
					c.JSON(http.StatusBadRequest, response)
					return
				}

				hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
				if err != nil {
					response := gin.H{"Message": "ERROR: BCRYPT ERROR"}
					c.JSON(http.StatusInternalServerError, response)
					return
				}

				teamID := c.MustGet("id").(uint)
				oldTeam := models.Team{Model: gorm.Model{ID: teamID}}
				newTeam := models.Team{HashedPassword: hashedPassword}
				if err := db.Where(&oldTeam).Updates(&newTeam).Error; err != nil {
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

func CompetitionRegistration() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.GetDB()
		role := c.MustGet("role").(middlewares.AuthRole)

		switch role {
		case middlewares.Team:
			{
				query := CompetitionRegistrationQuery{}
				if err := c.BindQuery(&query); err != nil {
					response := gin.H{"Message": "ERROR: BAD REQUEST"}
					c.JSON(http.StatusBadRequest, response)
					return
				}

				teamID := c.MustGet("id").(uint)
				oldTeam := models.Team{Model: gorm.Model{ID: teamID}}
				newTeam := models.Team{TeamCategory: query.TeamCategory}
				if err := db.Where(&oldTeam).Updates(&newTeam).Error; err != nil {
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

func ChangeStatusTeamHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.GetDB()
		role := c.MustGet("role").(middlewares.AuthRole)

		switch role {
		case middlewares.Admin:
			{
				request := ChangeStatusTeamRequest{}
				if err := c.BindJSON(&request); err != nil {
					response := gin.H{"Message": "ERROR: BAD REQUEST"}
					c.JSON(http.StatusBadRequest, response)
					return
				}

				query := ChangeStatusTeamQuery{}
				if err := c.BindQuery(&query); err != nil {
					response := gin.H{"Message": "ERROR: BAD REQUEST"}
					c.JSON(http.StatusBadRequest, response)
					return
				}

				adminID := c.MustGet("id").(uint)
				oldTeam := models.Team{Model: gorm.Model{ID: query.TeamID}}
				newTeam := models.Team{Status: request.Status, AdminID: adminID}
				if err := db.Where(&oldTeam).Updates(&newTeam).Error; err != nil {
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
