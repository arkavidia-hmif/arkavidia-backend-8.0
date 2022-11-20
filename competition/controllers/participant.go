package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"arkavidia-backend-8.0/competition/models"
	databaseService "arkavidia-backend-8.0/competition/services/database"
)

type AddMemberRequest struct {
	Name           string                             `json:"name"`
	Email          string                             `json:"email"`
	CareerInterest []models.ParticipantCareerInterest `json:"career_interest"`
	Role           models.MembershipRole              `json:"role"`
}

type ChangeCareerInterestRequest struct {
	ParticipantID  uint                               `json:"participant_id"`
	CareerInterest []models.ParticipantCareerInterest `json:"career_interest"`
}

type ChangeRoleRequest struct {
	ParticipantID uint                  `json:"participant_id"`
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
				participant := models.Participant{Model: gorm.Model{ID: membership.ParticipantID}}
				if err := tx.Find(&participant).Error; err != nil {
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
			participant = models.Participant{Name: request.Name, Email: request.Email, CareerInterest: request.CareerInterest}
			if err := tx.Find(&participant).Error; err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					return err
				}
			}
			if err := tx.Create(&participant).Error; err != nil {
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

		membership := models.Membership{TeamID: teamID, ParticipantID: request.ParticipantID}
		if err := db.Find(&membership).Error; err != nil {
			response := gin.H{"Message": "ERROR: BAD REQUEST"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		oldParticipant := models.Participant{Model: gorm.Model{ID: request.ParticipantID}}
		newParticipant := models.Participant{CareerInterest: request.CareerInterest}
		if err := db.Find(&oldParticipant).Updates(&newParticipant).Error; err != nil {
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

		oldMembership := models.Membership{TeamID: teamID, ParticipantID: request.ParticipantID}
		if err := db.Find(&oldMembership).Error; err != nil {
			response := gin.H{"Message": "ERROR: BAD REQUEST"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		newMembership := models.Membership{Role: request.Role}
		if err := db.Find(&oldMembership).Updates(&newMembership).Error; err != nil {
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

		membership := models.Membership{TeamID: teamID, ParticipantID: request.ParticipantID}
		if err := db.Find(&membership).Error; err != nil {
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
