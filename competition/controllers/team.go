package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	authConfig "arkavidia-backend-8.0/competition/config/authentication"
	"arkavidia-backend-8.0/competition/middlewares"
	"arkavidia-backend-8.0/competition/models"
	databaseService "arkavidia-backend-8.0/competition/services/database"
	"arkavidia-backend-8.0/competition/utils/mail"
	"arkavidia-backend-8.0/competition/utils/sanitizer"
)

type SignInTeamRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,ascii"`
}

type SignUpTeamRequest struct {
	Username string             `json:"username" binding:"required,alphanum"`
	Password string             `json:"password" binding:"required,ascii"`
	TeamName string             `json:"team_name" binding:"required,ascii"`
	Members  []SignUpMembership `json:"member_list" binding:"required,dive"`
}

type GetTeamQuery struct {
	TeamID uint `form:"team_id" field:"team_id" binding:"required,gt=0"`
}

type GetAllTeamsQuery struct {
	Page         int                 `form:"page" field:"page" binding:"required,gt=0"`
	Size         int                 `form:"size" field:"size" binding:"required,gt=0"`
	TeamCategory models.TeamCategory `form:"team_category" field:"team_category" binding:"required,oneof=competitive-programming datavidia uxvidia arkalogica"`
}

type ChangePasswordRequest struct {
	Password string `json:"password" binding:"required,ascii"`
}

type CompetitionRegistrationQuery struct {
	TeamCategory models.TeamCategory `form:"competition" field:"competition" binding:"required,oneof=competitive-programming datavidia uxvidia arkalogica"`
}

type ChangeStatusTeamQuery struct {
	TeamID uint `form:"team_id" field:"team_id" binding:"required,gt=0"`
}

type ChangeStatusTeamRequest struct {
	Status models.TeamStatus `json:"status" binding:"required,oneof=waiting-for-evaluation passed eliminated"`
}

func SignInTeamHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.DB.GetConnection()
		config := authConfig.Config.GetMetadata()
		response := sanitizer.Response[string]{}

		request := SignInTeamRequest{}
		if err := c.ShouldBindJSON(&request); err != nil {
			response.Message = "ERROR: BAD REQUEST"
			c.AbortWithStatusJSON(http.StatusBadRequest, sanitizer.SanitizeStruct(response))
			return
		}

		condition := models.Team{Username: request.Username}
		team := models.Team{}
		if err := db.Where(&condition).Find(&team).Error; err != nil {
			response.Message = "ERROR: INVALID USERNAME OR PASSWORD"
			c.AbortWithStatusJSON(http.StatusUnauthorized, sanitizer.SanitizeStruct(response))
			return
		}

		if err := bcrypt.CompareHashAndPassword(team.HashedPassword, []byte(request.Password)); err != nil {
			response.Message = "ERROR: INVALID USERNAME OR PASSWORD"
			c.AbortWithStatusJSON(http.StatusUnauthorized, sanitizer.SanitizeStruct(response))
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
			response.Message = "ERROR: JWT SIGNING ERROR"
			c.AbortWithStatusJSON(http.StatusInternalServerError, sanitizer.SanitizeStruct(response))
			return
		}

		response.Message = "SUCCESS"
		response.Data = signedAuthToken
		c.JSON(http.StatusCreated, sanitizer.SanitizeStruct(response))
	}
}

func SignUpTeamHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.DB.GetConnection()
		config := authConfig.Config.GetMetadata()
		response := sanitizer.Response[string]{}

		request := SignUpTeamRequest{}
		if err := c.ShouldBindJSON(&request); err != nil {
			response.Message = "ERROR: BAD REQUEST"
			c.AbortWithStatusJSON(http.StatusBadRequest, sanitizer.SanitizeStruct(response))
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			response.Message = "ERROR: BCRYPT ERROR"
			c.AbortWithStatusJSON(http.StatusInternalServerError, sanitizer.SanitizeStruct(response))
			return
		}

		// Validate Username
		conditionUsername := models.Team{Username: request.Username}
		team := models.Team{}
		if err := db.Where(&conditionUsername).Find(&team).Error; err != nil {
			response.Message = "ERROR: BAD REQUEST"
			c.AbortWithStatusJSON(http.StatusBadRequest, sanitizer.SanitizeStruct(response))
			return
		}
		if team.Username != "" {
			response.Message = "ERROR: USERNAME EXISTED"
			c.AbortWithStatusJSON(http.StatusBadRequest, sanitizer.SanitizeStruct(response))
			return
		}

		// Validate Team Name
		conditionTeamName := models.Team{TeamName: request.TeamName}
		team = models.Team{}
		if err := db.Where(&conditionTeamName).Find(&team).Error; err != nil {
			response.Message = "ERROR: BAD REQUEST"
			c.AbortWithStatusJSON(http.StatusBadRequest, sanitizer.SanitizeStruct(response))
			return
		}
		if team.TeamName != "" {
			response.Message = "ERROR: TEAMNAME EXISTED"
			c.AbortWithStatusJSON(http.StatusBadRequest, sanitizer.SanitizeStruct(response))
			return
		}

		team = models.Team{Username: request.Username, HashedPassword: hashedPassword, TeamName: request.TeamName, Status: models.WaitingForEvaluation}
		if err := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(&team).Error; err != nil {
				return err
			}

			for _, member := range request.Members {
				conditionParticipant := models.Participant{Name: member.Name, Email: member.Email, CareerInterest: member.CareerInterests, Status: models.WaitingForVerification}
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
			response.Message = "BAD REQUEST"
			c.AbortWithStatusJSON(http.StatusBadRequest, sanitizer.SanitizeStruct(response))
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
			response.Message = "ERROR: JWT SIGNING ERROR"
			c.AbortWithStatusJSON(http.StatusInternalServerError, sanitizer.SanitizeStruct(response))
			return
		}

		// Asynchronously mail to every member registered on the team
		for _, member := range request.Members {
			mail.Broker.AddMailToBroker(mail.MailParameters{Email: member.Email})
		}

		response.Message = "SUCCESS"
		response.Data = signedAuthToken
		c.JSON(http.StatusCreated, sanitizer.SanitizeStruct(response))
	}
}

func GetTeamHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.DB.GetConnection()
		response := sanitizer.Response[models.Team]{}

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
				query := GetTeamQuery{}
				if err := c.ShouldBindQuery(&query); err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, sanitizer.SanitizeStruct(response))
					return
				}

				condition := models.Team{Model: gorm.Model{ID: query.TeamID}}
				team := models.Team{}
				if err := db.Preload("Memberships").Where(&condition).Find(&team).Error; err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, sanitizer.SanitizeStruct(response))
					return
				}

				response.Message = "SUCCESS"
				response.Data = team
				c.JSON(http.StatusOK, sanitizer.SanitizeStruct(response))
				return
			}
		case middlewares.Team:
			{
				value, exists := c.Get("id")
				if !exists {
					response.Message = "UNAUTHORIZED"
					c.AbortWithStatusJSON(http.StatusUnauthorized, sanitizer.SanitizeStruct(response))
					return
				}

				teamID := value.(uint)
				condition := models.Team{Model: gorm.Model{ID: teamID}}
				team := models.Team{}
				if err := db.Preload("Memberships").Where(&condition).Find(&team).Error; err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, sanitizer.SanitizeStruct(response))
					return
				}

				response.Message = "SUCCESS"
				response.Data = team
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

func GetAllTeamsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.DB.GetConnection()
		response := sanitizer.Response[[]models.Team]{}

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
				query := GetAllTeamsQuery{}
				if err := c.ShouldBindQuery(&query); err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, sanitizer.SanitizeStruct(response))
					return
				}

				offset := (query.Page - 1) * query.Size
				limit := query.Size
				condition := models.Team{TeamCategory: query.TeamCategory}
				teams := []models.Team{}
				if err := db.Where(&condition).Offset(offset).Limit(limit).Find(&teams).Error; err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, sanitizer.SanitizeStruct(response))
					return
				}

				response.Message = "SUCCESS"
				response.Data = teams
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

func ChangePasswordHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.DB.GetConnection()
		response := sanitizer.Response[models.Team]{}

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
				request := ChangePasswordRequest{}
				if err := c.ShouldBindJSON(&request); err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, sanitizer.SanitizeStruct(response))
					return
				}

				hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
				if err != nil {
					response.Message = "ERROR: BCRYPT ERROR"
					c.AbortWithStatusJSON(http.StatusInternalServerError, sanitizer.SanitizeStruct(response))
					return
				}

				value, exists := c.Get("id")
				if !exists {
					response.Message = "UNAUTHORIZED"
					c.AbortWithStatusJSON(http.StatusUnauthorized, sanitizer.SanitizeStruct(response))
					return
				}

				teamID := value.(uint)
				oldTeam := models.Team{Model: gorm.Model{ID: teamID}}
				newTeam := models.Team{HashedPassword: hashedPassword}
				if err := db.Where(&oldTeam).Updates(&newTeam).Error; err != nil {
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

func CompetitionRegistration() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.DB.GetConnection()
		response := sanitizer.Response[models.Team]{}

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
				query := CompetitionRegistrationQuery{}
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
				oldTeam := models.Team{Model: gorm.Model{ID: teamID}}
				newTeam := models.Team{TeamCategory: query.TeamCategory}
				if err := db.Where(&oldTeam).Updates(&newTeam).Error; err != nil {
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

func ChangeStatusTeamHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.DB.GetConnection()
		response := sanitizer.Response[models.Team]{}

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
				request := ChangeStatusTeamRequest{}
				if err := c.ShouldBindJSON(&request); err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, sanitizer.SanitizeStruct(response))
					return
				}

				query := ChangeStatusTeamQuery{}
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
				oldTeam := models.Team{Model: gorm.Model{ID: query.TeamID}}
				newTeam := models.Team{Status: request.Status, AdminID: adminID}
				if err := db.Where(&oldTeam).Updates(&newTeam).Error; err != nil {
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
