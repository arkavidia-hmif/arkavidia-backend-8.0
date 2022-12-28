package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"arkavidia-backend-8.0/competition/middlewares"
	"arkavidia-backend-8.0/competition/models"
	databaseService "arkavidia-backend-8.0/competition/services/database"
	"arkavidia-backend-8.0/competition/utils/sanitizer"
)

type GetMemberQuery struct {
	TeamID uint `form:"team_id" field:"team_id" binding:"required,gt=0"`
}

type GetAllMembersQuery struct {
	Page int `form:"page" field:"page" binding:"required,gt=0"`
	Size int `form:"size" field:"size" binding:"required,gt=0"`
}

type AddMemberRequest struct {
	Name            string                            `json:"name" binding:"required,ascii"`
	Email           string                            `json:"email" binding:"required,email"`
	CareerInterests models.ParticipantCareerInterests `json:"career_interest" binding:"required,dive,oneof=software-engineering product-management ui-designer ux-designer ux-researcher it-consultant game-developer cyber-security business-analyst business-intelligence data-scientist data-analyst"`
	Role            models.MembershipRole             `json:"role" binding:"required,oneof=leader member"`
}

type ChangeCareerInterestQuery struct {
	ParticipantID uint `form:"participant_id" field:"participant_id" binding:"required,gt=0"`
}

type ChangeCareerInterestRequest struct {
	CareerInterests models.ParticipantCareerInterests `json:"career_interest" binding:"required,dive,oneof=software-engineering product-management ui-designer ux-designer ux-researcher it-consultant game-developer cyber-security business-analyst business-intelligence data-scientist data-analyst"`
}

type ChangeRoleQuery struct {
	ParticipantID uint `form:"participant_id" field:"participant_id" binding:"required,gt=0"`
}

type ChangeRoleRequest struct {
	Role models.MembershipRole `json:"role" binding:"required,oneof=leader member"`
}

type ChangeStatusParticipantQuery struct {
	ParticipantID uint `form:"participant_id" field:"participant_id" binding:"required,gt=0"`
}

type ChangeStatusParticipantRequest struct {
	Status models.ParticipantStatus `json:"status" binding:"required,oneof=waiting-for-verification verified declined"`
}

type DeleteMemberRequest struct {
	ParticipantID uint `json:"participant_id" binding:"required,gt=0"`
}

func GetMemberHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.DB.GetConnection()
		response := sanitizer.Response[[]models.Participant]{}

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
				query := GetMemberQuery{}
				if err := c.ShouldBindQuery(&query); err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, sanitizer.SanitizeStruct(response))
					return
				}

				condition := models.Membership{TeamID: query.TeamID}
				memberships := []models.Membership{}
				if err := db.Where(&condition).Find(&memberships).Error; err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, sanitizer.SanitizeStruct(response))
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
					c.AbortWithStatusJSON(http.StatusBadRequest, sanitizer.SanitizeStruct(response))
					return
				}

				response.Message = "SUCCESS"
				response.Data = participants
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
				condition := models.Membership{TeamID: teamID}
				memberships := []models.Membership{}
				if err := db.Where(&condition).Find(&memberships).Error; err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, sanitizer.SanitizeStruct(response))
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
					c.AbortWithStatusJSON(http.StatusBadRequest, sanitizer.SanitizeStruct(response))
					return
				}

				response.Message = "SUCCESS"
				response.Data = participants
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

func GetAllMembersHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.DB.GetConnection()
		response := sanitizer.Response[[]models.Participant]{}

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
				query := GetAllMembersQuery{}
				if err := c.ShouldBindQuery(&query); err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, sanitizer.SanitizeStruct(response))
					return
				}

				offset := (query.Page - 1) * query.Size
				limit := query.Size
				participants := []models.Participant{}

				if err := db.Offset(offset).Limit(limit).Find(&participants).Error; err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, sanitizer.SanitizeStruct(response))
					return
				}

				response.Message = "SUCCESS"
				response.Data = participants
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

func AddMemberHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.DB.GetConnection()
		response := sanitizer.Response[models.Participant]{}

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
				request := AddMemberRequest{}
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
				participant := models.Participant{}
				membership := models.Membership{}
				if err := db.Transaction(func(tx *gorm.DB) error {
					condition := models.Participant{Name: request.Name, Email: request.Email, CareerInterest: request.CareerInterests, Status: models.WaitingForVerification}
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
					c.AbortWithStatusJSON(http.StatusBadRequest, sanitizer.SanitizeStruct(response))
					return
				}

				response.Message = "SUCCESS"
				response.Data = participant
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

func ChangeCareerInterestHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.DB.GetConnection()
		response := sanitizer.Response[models.Participant]{}

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
				request := ChangeCareerInterestRequest{}
				if err := c.ShouldBindJSON(&request); err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, sanitizer.SanitizeStruct(response))
					return
				}

				query := ChangeCareerInterestQuery{}
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
				condition := models.Membership{TeamID: teamID, ParticipantID: query.ParticipantID}
				membership := models.Membership{}
				if err := db.Where(&condition).Find(&membership).Error; err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, sanitizer.SanitizeStruct(response))
					return
				}

				oldParticipant := models.Participant{Model: gorm.Model{ID: query.ParticipantID}}
				newParticipant := models.Participant{CareerInterest: request.CareerInterests}
				if err := db.Where(&oldParticipant).Updates(&newParticipant).Error; err != nil {
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

func ChangeRoleHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.DB.GetConnection()
		response := sanitizer.Response[models.Participant]{}

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
				request := ChangeRoleRequest{}
				if err := c.ShouldBindJSON(&request); err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, sanitizer.SanitizeStruct(response))
					return
				}

				query := ChangeRoleQuery{}
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
				oldMembership := models.Membership{TeamID: teamID, ParticipantID: query.ParticipantID}
				newMembership := models.Membership{Role: request.Role}
				if err := db.Where(&oldMembership).Updates(&newMembership).Error; err != nil {
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

func ChangeStatusParticipantHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.DB.GetConnection()
		response := sanitizer.Response[models.Participant]{}

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
				request := ChangeStatusParticipantRequest{}
				if err := c.ShouldBindJSON(&request); err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, sanitizer.SanitizeStruct(response))
					return
				}

				query := ChangeStatusParticipantQuery{}
				if err := c.ShouldBindQuery(&query); err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, sanitizer.SanitizeStruct(response))
					return
				}

				oldParticipant := models.Participant{Model: gorm.Model{ID: query.ParticipantID}}
				newParticipant := models.Participant{Status: request.Status}
				if err := db.Where(&oldParticipant).Updates(&newParticipant).Error; err != nil {
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

func DeleteParticipantHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.DB.GetConnection()
		response := sanitizer.Response[models.Participant]{}

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
				request := DeleteMemberRequest{}
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
				condition := models.Membership{TeamID: teamID, ParticipantID: request.ParticipantID}
				membership := models.Membership{}
				if err := db.Where(&condition).Find(&membership).Error; err != nil {
					response.Message = "ERROR: BAD REQUEST"
					c.AbortWithStatusJSON(http.StatusBadRequest, sanitizer.SanitizeStruct(response))
					return
				}

				if err := db.Delete(&membership); err != nil {
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
