package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"gorm.io/gorm"

	"arkavidia-backend-8.0/competition/middlewares"
	"arkavidia-backend-8.0/competition/models"
	databaseService "arkavidia-backend-8.0/competition/services/database"
)

type GetMemberQuery struct {
	TeamID uint `form:"team_id" field:"team_id" binding:"required,gt=0"`
}

type GetAllMembersQuery struct {
	Page int `form:"page" field:"page" binding:"required,gt=0"`
	Size int `form:"size" field:"size" binding:"required,gt=0"`
}

type AddMemberRequest struct {
	Name            string                             `json:"name" binding:"required,ascii"`
	Email           string                             `json:"email" binding:"required,email"`
	CareerInterests []models.ParticipantCareerInterest `json:"career_interest" binding:"required,dive,oneof=software-engineering product-management ui-designer ux-designer ux-researcher it-consultant game-developer cyber-security business-analyst business-intelligence data-scientist data-analyst"`
	Role            models.MembershipRole              `json:"role" binding:"required,oneof=leader member"`
}

type ChangeCareerInterestQuery struct {
	ParticipantID uint `form:"participant_id" field:"participant_id" binding:"required,gt=0"`
}

type ChangeCareerInterestRequest struct {
	CareerInterests []models.ParticipantCareerInterest `json:"career_interest" binding:"required,dive,oneof=software-engineering product-management ui-designer ux-designer ux-researcher it-consultant game-developer cyber-security business-analyst business-intelligence data-scientist data-analyst"`
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
		role := c.MustGet("role").(middlewares.AuthRole)

		switch role {
		case middlewares.Admin:
			{
				query := GetMemberQuery{}
				if err := c.BindQuery(&query); err != nil {
					response := gin.H{"Message": "ERROR: BAD REQUEST"}
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				condition := models.Membership{TeamID: query.TeamID}
				memberships := []models.Membership{}
				if err := db.Where(&condition).Find(&memberships).Error; err != nil {
					response := gin.H{"Message": "ERROR: BAD REQUEST"}
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
					response := gin.H{"Message": "ERROR: BAD REQUEST"}
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				response := gin.H{"Message": "SUCCESS", "Data": participants}
				c.JSON(http.StatusOK, response)
				return
			}
		case middlewares.Team:
			{
				teamID := c.MustGet("id").(uint)
				condition := models.Membership{TeamID: teamID}
				memberships := []models.Membership{}
				if err := db.Where(&condition).Find(&memberships).Error; err != nil {
					response := gin.H{"Message": "ERROR: BAD REQUEST"}
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
					response := gin.H{"Message": "ERROR: BAD REQUEST"}
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				response := gin.H{"Message": "SUCCESS", "Data": participants}
				c.JSON(http.StatusOK, response)
				return
			}
		default:
			{
				response := gin.H{"Message": "ERROR: INVALID ROLE"}
				c.AbortWithStatusJSON(http.StatusUnauthorized, response)
				return
			}
		}
	}
}

func GetAllMembersHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.DB.GetConnection()
		role := c.MustGet("role").(middlewares.AuthRole)

		switch role {
		case middlewares.Admin:
			{
				query := GetAllMembersQuery{}
				if err := c.BindQuery(&query); err != nil {
					response := gin.H{"Message": "ERROR: BAD REQUEST"}
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				offset := (query.Page - 1) * query.Size
				limit := query.Size
				members := []models.Membership{}

				if err := db.Offset(offset).Limit(limit).Find(&members).Error; err != nil {
					response := gin.H{"Message": "ERROR: BAD REQUEST"}
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				response := gin.H{"Message": "SUCCESS", "Data": members}
				c.JSON(http.StatusOK, response)
				return
			}
		default:
			{
				response := gin.H{"Message": "ERROR: INVALID ROLE"}
				c.AbortWithStatusJSON(http.StatusUnauthorized, response)
				return
			}
		}
	}
}

func AddMemberHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.DB.GetConnection()
		role := c.MustGet("role").(middlewares.AuthRole)

		switch role {
		case middlewares.Team:
			{
				request := AddMemberRequest{}
				if err := c.ShouldBindJSON(&request); err != nil {
					response := gin.H{"Message": "ERROR: BAD REQUEST"}
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				participant := models.Participant{}
				membership := models.Membership{}
				if err := db.Transaction(func(tx *gorm.DB) error {
					var pqStringArray pq.StringArray
					for _, CareerInterest := range request.CareerInterests {
						pqStringArray = append(pqStringArray, string(CareerInterest))
					}

					condition := models.Participant{Name: request.Name, Email: request.Email, CareerInterest: pqStringArray, Status: models.WaitingForVerification}
					if err := tx.FirstOrCreate(&participant, &condition).Error; err != nil {
						return err
					}

					teamID := c.MustGet("id").(uint)
					membership = models.Membership{TeamID: teamID, ParticipantID: participant.ID, Role: request.Role}
					if err := tx.Create(&membership).Error; err != nil {
						return err
					}

					participant.Memberships = append(participant.Memberships, membership)

					return nil

				}); err != nil {
					response := gin.H{"Message": "ERROR: BAD REQUEST"}
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				response := gin.H{"Message": "SUCCESS", "Data": participant}
				c.JSON(http.StatusCreated, response)
				return
			}
		default:
			{
				response := gin.H{"Message": "ERROR: INVALID ROLE"}
				c.AbortWithStatusJSON(http.StatusUnauthorized, response)
				return
			}
		}
	}
}

func ChangeCareerInterestHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.DB.GetConnection()
		role := c.MustGet("role").(middlewares.AuthRole)

		switch role {
		case middlewares.Team:
			{
				request := ChangeCareerInterestRequest{}
				if err := c.ShouldBindJSON(&request); err != nil {
					response := gin.H{"Message": "ERROR: BAD REQUEST"}
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				query := ChangeCareerInterestQuery{}
				if err := c.BindQuery(&query); err != nil {
					response := gin.H{"Message": "ERROR: BAD REQUEST"}
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				teamID := c.MustGet("id").(uint)
				condition := models.Membership{TeamID: teamID, ParticipantID: query.ParticipantID}
				membership := models.Membership{}
				if err := db.Where(&condition).Find(&membership).Error; err != nil {
					response := gin.H{"Message": "ERROR: BAD REQUEST"}
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				var pqStringArray pq.StringArray
				for _, CareerInterest := range request.CareerInterests {
					pqStringArray = append(pqStringArray, string(CareerInterest))
				}

				oldParticipant := models.Participant{Model: gorm.Model{ID: query.ParticipantID}}
				newParticipant := models.Participant{CareerInterest: pqStringArray}
				if err := db.Where(&oldParticipant).Updates(&newParticipant).Error; err != nil {
					response := gin.H{"Message": "ERROR: BAD REQUEST"}
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				response := gin.H{"Message": "SUCCESS"}
				c.JSON(http.StatusOK, response)
				return
			}
		default:
			{
				response := gin.H{"Message": "ERROR: INVALID ROLE"}
				c.AbortWithStatusJSON(http.StatusUnauthorized, response)
				return
			}
		}
	}
}

func ChangeRoleHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.DB.GetConnection()
		role := c.MustGet("role").(middlewares.AuthRole)

		switch role {
		case middlewares.Team:
			{
				request := ChangeRoleRequest{}
				if err := c.ShouldBindJSON(&request); err != nil {
					response := gin.H{"Message": "ERROR: BAD REQUEST"}
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				query := ChangeRoleQuery{}
				if err := c.BindQuery(&query); err != nil {
					response := gin.H{"Message": "ERROR: BAD REQUEST"}
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				teamID := c.MustGet("id").(uint)
				oldMembership := models.Membership{TeamID: teamID, ParticipantID: query.ParticipantID}
				newMembership := models.Membership{Role: request.Role}
				if err := db.Where(&oldMembership).Updates(&newMembership).Error; err != nil {
					response := gin.H{"Message": "ERROR: BAD REQUEST"}
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				response := gin.H{"Message": "SUCCESS"}
				c.JSON(http.StatusOK, response)
				return
			}
		default:
			{
				response := gin.H{"Message": "ERROR: INVALID ROLE"}
				c.AbortWithStatusJSON(http.StatusUnauthorized, response)
				return
			}
		}
	}
}

func ChangeStatusParticipantHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.DB.GetConnection()
		role := c.MustGet("role").(middlewares.AuthRole)

		switch role {
		case middlewares.Team:
			{
				request := ChangeStatusParticipantRequest{}
				if err := c.ShouldBindJSON(&request); err != nil {
					response := gin.H{"Message": "ERROR: BAD REQUEST"}
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				query := ChangeStatusParticipantQuery{}
				if err := c.BindQuery(&query); err != nil {
					response := gin.H{"Message": "ERROR: BAD REQUEST"}
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				oldParticipant := models.Participant{Model: gorm.Model{ID: query.ParticipantID}}
				newParticipant := models.Participant{Status: request.Status}
				if err := db.Where(&oldParticipant).Updates(&newParticipant).Error; err != nil {
					response := gin.H{"Message": "ERROR: BAD REQUEST"}
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				response := gin.H{"Message": "SUCCESS"}
				c.JSON(http.StatusOK, response)
				return
			}
		default:
			{
				response := gin.H{"Message": "ERROR: INVALID ROLE"}
				c.AbortWithStatusJSON(http.StatusUnauthorized, response)
				return
			}
		}
	}
}

func DeleteParticipantHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.DB.GetConnection()
		role := c.MustGet("role").(middlewares.AuthRole)

		switch role {
		case middlewares.Team:
			{
				request := DeleteMemberRequest{}
				if err := c.ShouldBindJSON(&request); err != nil {
					response := gin.H{"Message": "ERROR: BAD REQUEST"}
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				teamID := c.MustGet("id").(uint)
				condition := models.Membership{TeamID: teamID, ParticipantID: request.ParticipantID}
				membership := models.Membership{}
				if err := db.Where(&condition).Find(&membership).Error; err != nil {
					response := gin.H{"Message": "ERROR: BAD REQUEST"}
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				if err := db.Delete(&membership); err != nil {
					response := gin.H{"Message": "ERROR: BAD REQUEST"}
					c.AbortWithStatusJSON(http.StatusBadRequest, response)
					return
				}

				response := gin.H{"Message": "SUCCESS"}
				c.JSON(http.StatusOK, response)
				return
			}
		default:
			{
				response := gin.H{"Message": "ERROR: INVALID ROLE"}
				c.AbortWithStatusJSON(http.StatusUnauthorized, response)
				return
			}
		}
	}
}
