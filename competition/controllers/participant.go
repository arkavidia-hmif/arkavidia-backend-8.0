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
