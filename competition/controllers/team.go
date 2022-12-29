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
	"arkavidia-backend-8.0/competition/repository"
	databaseService "arkavidia-backend-8.0/competition/services/database"
	"arkavidia-backend-8.0/competition/types"
	"arkavidia-backend-8.0/competition/utils/mail"
)

func SignInTeamHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.DB.GetConnection()
		config := authConfig.Config.GetMetadata()
		response := repository.Response[string]{}

		request := repository.SignInTeamRequest{}
		if err := c.ShouldBindJSON(&request); err != nil {
			response.Message = "ERROR: BAD REQUEST"
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		condition := models.Team{Username: request.Username}
		team := models.Team{}
		if err := db.Where(&condition).Find(&team).Error; err != nil {
			response.Message = "ERROR: INVALID USERNAME OR PASSWORD"
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		if err := bcrypt.CompareHashAndPassword(team.HashedPassword, []byte(request.Password)); err != nil {
			response.Message = "ERROR: INVALID USERNAME OR PASSWORD"
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
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
			c.AbortWithStatusJSON(http.StatusInternalServerError, response)
			return
		}

		response.Message = "SUCCESS"
		response.Data = signedAuthToken
		c.JSON(http.StatusCreated, response)
	}
}

func SignUpTeamHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.DB.GetConnection()
		config := authConfig.Config.GetMetadata()
		response := repository.Response[string]{}

		request := repository.SignUpTeamRequest{}
		if err := c.ShouldBindJSON(&request); err != nil {
			response.Message = "ERROR: BAD REQUEST"
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		// Validate Username
		conditionUsername := models.Team{Username: request.Username}
		team := models.Team{}
		if err := db.Where(&conditionUsername).Find(&team).Error; err != nil {
			response.Message = "ERROR: BAD REQUEST"
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}
		if team.Username != "" {
			response.Message = "ERROR: USERNAME EXISTED"
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		// Validate Team Name
		conditionTeamName := models.Team{TeamName: request.TeamName}
		team = models.Team{}
		if err := db.Where(&conditionTeamName).Find(&team).Error; err != nil {
			response.Message = "ERROR: BAD REQUEST"
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}
		if team.TeamName != "" {
			response.Message = "ERROR: TEAMNAME EXISTED"
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		encryptedString := []byte(request.Password)
		team = models.Team{Username: request.Username, HashedPassword: encryptedString, TeamName: request.TeamName, Status: types.WaitingForEvaluation}
		if err := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(&team).Error; err != nil {
				return err
			}

			for _, member := range request.Members {
				conditionParticipant := models.Participant{Name: member.Name, Email: member.Email, CareerInterest: member.CareerInterests, Status: types.WaitingForVerification}
				participant := models.Participant{}
				if err := tx.FirstOrCreate(&participant, &conditionParticipant).Error; err != nil {
					return err
				}

				membership := models.Membership{TeamID: team.ID, ParticipantID: participant.ID, Role: member.Role}
				if err := tx.Create(&membership).Error; err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			response.Message = "BAD REQUEST"
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
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
			c.AbortWithStatusJSON(http.StatusInternalServerError, response)
			return
		}

		// Asynchronously mail to every member registered on the team
		for _, member := range request.Members {
			mail.Broker.AddMailToBroker(mail.MailParameters{Email: member.Email})
		}

		response.Message = "SUCCESS"
		response.Data = signedAuthToken
		c.JSON(http.StatusCreated, response)
	}
}

func GetTeamHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.DB.GetConnection()
		response := repository.Response[models.Team]{}

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
				query := repository.GetTeamQuery{}
				if err := c.ShouldBindQuery(&query); err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				condition := models.Team{Model: gorm.Model{ID: query.TeamID}}
				team := models.Team{}
				if err := db.Preload("Memberships").Where(&condition).Find(&team).Error; err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				response.Message = "SUCCESS"
				response.Data = team
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
				condition := models.Team{Model: gorm.Model{ID: teamID}}
				team := models.Team{}
				if err := db.Preload("Memberships").Where(&condition).Find(&team).Error; err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				response.Message = "SUCCESS"
				response.Data = team
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

func GetAllTeamsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.DB.GetConnection()
		response := repository.Response[[]models.Team]{}

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
				query := repository.GetAllTeamsQuery{}
				if err := c.ShouldBindQuery(&query); err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				offset := (query.Page - 1) * query.Size
				limit := query.Size
				condition := models.Team{TeamCategory: query.TeamCategory}
				teams := []models.Team{}
				if err := db.Where(&condition).Offset(offset).Limit(limit).Find(&teams).Error; err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				response.Message = "SUCCESS"
				response.Data = teams
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

func ChangePasswordHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.DB.GetConnection()
		response := repository.Response[models.Team]{}

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
				request := repository.ChangePasswordRequest{}
				if err := c.ShouldBindJSON(&request); err != nil {
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
				oldTeam := models.Team{Model: gorm.Model{ID: teamID}}
				newTeam := models.Team{HashedPassword: []byte(request.Password)}
				if err := db.Where(&oldTeam).Updates(&newTeam).Error; err != nil {
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

func CompetitionRegistration() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.DB.GetConnection()
		response := repository.Response[models.Team]{}

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
				query := repository.CompetitionRegistrationQuery{}
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
				oldTeam := models.Team{Model: gorm.Model{ID: teamID}}
				newTeam := models.Team{TeamCategory: query.TeamCategory}
				if err := db.Where(&oldTeam).Updates(&newTeam).Error; err != nil {
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

func ChangeStatusTeamHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.DB.GetConnection()
		response := repository.Response[models.Team]{}

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
				request := repository.ChangeStatusTeamRequest{}
				if err := c.ShouldBindJSON(&request); err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				query := repository.ChangeStatusTeamQuery{}
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

				adminID := value.(uint)
				oldTeam := models.Team{Model: gorm.Model{ID: query.TeamID}}
				newTeam := models.Team{Status: request.Status, AdminID: adminID}
				if err := db.Where(&oldTeam).Updates(&newTeam).Error; err != nil {
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
