package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"arkavidia-backend-8.0/competition/middlewares"
	"arkavidia-backend-8.0/competition/models"
	"arkavidia-backend-8.0/competition/repository"
	databaseService "arkavidia-backend-8.0/competition/services/database"
	"arkavidia-backend-8.0/competition/types"
)

func GetMemberHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.DB.GetConnection()
		response := repository.Response[[]models.Participant]{}

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
				query := repository.GetMemberQuery{}
				if err := c.ShouldBindQuery(&query); err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				condition := models.Membership{TeamID: query.TeamID}
				memberships := []models.Membership{}
				if err := db.Where(&condition).Find(&memberships).Error; err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				participants := []models.Participant{}
				if err := db.Transaction(func(tx *gorm.DB) error {
					for _, membership := range memberships {
						condition := models.Participant{Model: gorm.Model{ID: membership.ParticipantID}}
						participant := models.Participant{}
						if err := tx.Where(&condition).Find(&participant).Error; err != nil {
							return err
						}
						participant.Memberships = append(participant.Memberships, membership)
						participants = append(participants, participant)
					}
					return nil
				}); err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				response.Message = "SUCCESS"
				response.Data = participants
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
				condition := models.Membership{TeamID: teamID}
				memberships := []models.Membership{}
				if err := db.Where(&condition).Find(&memberships).Error; err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				participants := []models.Participant{}
				if err := db.Transaction(func(tx *gorm.DB) error {
					for _, membership := range memberships {
						condition := models.Participant{Model: gorm.Model{ID: membership.ParticipantID}}
						participant := models.Participant{}
						if err := tx.Where(&condition).Find(&participant).Error; err != nil {
							return err
						}
						participant.Memberships = append(participant.Memberships, membership)
						participants = append(participants, participant)
					}
					return nil
				}); err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				response.Message = "SUCCESS"
				response.Data = participants
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

func GetAllMembersHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.DB.GetConnection()
		response := repository.Response[[]models.Participant]{}

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
				query := repository.GetAllMembersQuery{}
				if err := c.ShouldBindQuery(&query); err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				offset := (query.Page - 1) * query.Size
				limit := query.Size
				participants := []models.Participant{}

				if err := db.Offset(offset).Limit(limit).Find(&participants).Error; err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				response.Message = "SUCCESS"
				response.Data = participants
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

func AddMemberHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.DB.GetConnection()
		response := repository.Response[models.Participant]{}

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
				request := repository.AddMemberRequest{}
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
				participant := models.Participant{}
				membership := models.Membership{}
				if err := db.Transaction(func(tx *gorm.DB) error {
					condition := models.Participant{Name: request.Name, Email: request.Email, CareerInterest: request.CareerInterests, Status: types.WaitingForVerification}
					if err := tx.FirstOrCreate(&participant, &condition).Error; err != nil {
						return err
					}

					membership = models.Membership{TeamID: teamID, ParticipantID: participant.ID, Role: request.Role}
					if err := tx.Create(&membership).Error; err != nil {
						return err
					}

					participant.Memberships = append(participant.Memberships, membership)

					return nil

				}); err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				response.Message = "SUCCESS"
				response.Data = participant
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

func ChangeCareerInterestHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.DB.GetConnection()
		response := repository.Response[models.Participant]{}

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
				request := repository.ChangeCareerInterestRequest{}
				if err := c.ShouldBindJSON(&request); err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				query := repository.ChangeCareerInterestQuery{}
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
				condition := models.Membership{TeamID: teamID, ParticipantID: query.ParticipantID}
				membership := models.Membership{}
				if err := db.Where(&condition).Find(&membership).Error; err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				oldParticipant := models.Participant{Model: gorm.Model{ID: query.ParticipantID}}
				newParticipant := models.Participant{CareerInterest: request.CareerInterests}
				if err := db.Where(&oldParticipant).Updates(&newParticipant).Error; err != nil {
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

func ChangeRoleHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.DB.GetConnection()
		response := repository.Response[models.Participant]{}

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
				request := repository.ChangeRoleRequest{}
				if err := c.ShouldBindJSON(&request); err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				query := repository.ChangeRoleQuery{}
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
				oldMembership := models.Membership{TeamID: teamID, ParticipantID: query.ParticipantID}
				newMembership := models.Membership{Role: request.Role}
				if err := db.Where(&oldMembership).Updates(&newMembership).Error; err != nil {
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

func ChangeStatusParticipantHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.DB.GetConnection()
		response := repository.Response[models.Participant]{}

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
				request := repository.ChangeStatusParticipantRequest{}
				if err := c.ShouldBindJSON(&request); err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				query := repository.ChangeStatusParticipantQuery{}
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

				condition := models.Membership{TeamID: teamID, ParticipantID: query.ParticipantID}
				membership := models.Membership{}
				if err := db.Where(&condition).Find(&membership).Error; err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				oldParticipant := models.Participant{Model: gorm.Model{ID: query.ParticipantID}}
				newParticipant := models.Participant{Status: request.Status}
				if err := db.Where(&oldParticipant).Updates(&newParticipant).Error; err != nil {
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

func DeleteParticipantHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.DB.GetConnection()
		response := repository.Response[models.Participant]{}

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
				request := repository.DeleteMemberRequest{}
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
				condition := models.Membership{TeamID: teamID, ParticipantID: request.ParticipantID}
				membership := models.Membership{}
				if err := db.Where(&condition).Find(&membership).Error; err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				if err := db.Delete(&membership); err != nil {
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
