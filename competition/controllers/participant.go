package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"gorm.io/gorm"

	"arkavidia-backend-8.0/competition/models"
	databaseService "arkavidia-backend-8.0/competition/services/database"
)

type AddMemberRequest struct {
	Name           string                `json:"name" binding:"required"`
	Email          string                `json:"email" binding:"required"`
	CareerInterest pq.StringArray        `json:"career_interest" binding:"required"`
	Role           models.MembershipRole `json:"role" binding:"required"`
}

type ChangeCareerInterestQuery struct {
	ParticipantID uint `form:"participant_id" field:"participant_id" binding:"required"`
}

type ChangeCareerInterestRequest struct {
	CareerInterest pq.StringArray `json:"career_interest" binding:"required"`
}

type ChangeRoleQuery struct {
	ParticipantID uint `form:"participant_id" field:"participant_id" binding:"required"`
}

type ChangeRoleRequest struct {
	Role          models.MembershipRole `json:"role"`
}

type DeleteMemberRequest struct {
	ParticipantID uint `json:"participant_id"`
}

func GetMemberHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.GetDB()
		teamID := c.MustGet("team_id").(uint)

		condition := models.Membership{TeamID: teamID}
		memberships := []models.Membership{}
		if err := db.Where(&condition).Find(&memberships).Error; err != nil {
			response := gin.H{"Message": "ERROR: BAD REQUEST"}
			c.JSON(http.StatusBadRequest, response)
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
			c.JSON(http.StatusBadRequest, response)
			return
		}

		response := gin.H{"Message": "SUCCESS", "Data": participants}
		c.JSON(http.StatusOK, response)
	}
}

func AddMemberHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.GetDB()
		teamID := c.MustGet("team_id").(uint)

		request := AddMemberRequest{}
		if err := c.BindJSON(&request); err != nil {
			response := gin.H{"Message": "ERROR: BAD REQUEST"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		participant := models.Participant{}
		membership := models.Membership{}
		if err := db.Transaction(func(tx *gorm.DB) error {
			condition := models.Participant{Name: request.Name, Email: request.Email, CareerInterest: request.CareerInterest}
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
			response := gin.H{"Message": "ERROR: BAD REQUEST"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		response := gin.H{"Message": "SUCCESS", "Data": participant}
		c.JSON(http.StatusCreated, response)
	}
}

func ChangeCareerInterestHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.GetDB()
		teamID := c.MustGet("team_id").(uint)

		request := ChangeCareerInterestRequest{}
		if err := c.BindJSON(&request); err != nil {
			response := gin.H{"Message": "ERROR: BAD REQUEST"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		query := ChangeCareerInterestQuery{}
		if err := c.BindQuery(&query); err != nil {
			response := gin.H{"Message": "ERROR: BAD REQUEST"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		condition := models.Membership{TeamID: teamID, ParticipantID: query.ParticipantID}
		membership := models.Membership{}
		if err := db.Where(&condition).Find(&membership).Error; err != nil {
			response := gin.H{"Message": "ERROR: BAD REQUEST"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		oldParticipant := models.Participant{Model: gorm.Model{ID: query.ParticipantID}}
		newParticipant := models.Participant{CareerInterest: request.CareerInterest}
		if err := db.Where(&oldParticipant).Updates(&newParticipant).Error; err != nil {
			response := gin.H{"Message": "ERROR: BAD REQUEST"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		response := gin.H{"Message": "SUCCESS"}
		c.JSON(http.StatusOK, response)
	}
}

func ChangeRoleInterestHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.GetDB()
		teamID := c.MustGet("team_id").(uint)

		request := ChangeRoleRequest{}
		if err := c.BindJSON(&request); err != nil {
			response := gin.H{"Message": "ERROR: BAD REQUEST"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		query := ChangeRoleQuery{}
		if err := c.BindQuery(&query); err != nil {
			response := gin.H{"Message": "ERROR: BAD REQUEST"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		oldMembership := models.Membership{TeamID: teamID, ParticipantID: query.ParticipantID}
		newMembership := models.Membership{Role: request.Role}
		if err := db.Where(&oldMembership).Updates(&newMembership).Error; err != nil {
			response := gin.H{"Message": "ERROR: BAD REQUEST"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		response := gin.H{"Message": "SUCCESS"}
		c.JSON(http.StatusOK, response)
	}
}

func DeleteParticipantHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.GetDB()
		teamID := c.MustGet("team_id").(uint)

		request := DeleteMemberRequest{}
		if err := c.BindJSON(&request); err != nil {
			response := gin.H{"Message": "ERROR: BAD REQUEST"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		condition := models.Membership{TeamID: teamID, ParticipantID: request.ParticipantID}
		membership := models.Membership{}
		if err := db.Where(&condition).Find(&membership).Error; err != nil {
			response := gin.H{"Message": "ERROR: BAD REQUEST"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		if err := db.Delete(&membership); err != nil {
			response := gin.H{"Message": "ERROR: BAD REQUEST"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		response := gin.H{"Message": "SUCCESS"}
		c.JSON(http.StatusOK, response)
	}
}
