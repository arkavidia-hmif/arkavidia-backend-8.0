package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"arkavidia-backend-8.0/competition/models"
	databaseService "arkavidia-backend-8.0/competition/services/database"
)

type AddParticipantRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ChangeCareerInterestRequest struct {
	ParticipantID  uint                               `json:"participant_id"`
	CareerInterest []models.ParticipantCareerInterest `json:"career_interest"`
}

func GetParticipantHandler() gin.HandlerFunc {
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

func AddParticipantHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.GetDB()

		request := AddParticipantRequest{}
		if err := c.BindJSON(&request); err != nil {
			response := gin.H{"Message": "ERROR: BAD REQUEST"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		participant := models.Participant{Name: request.Name, Email: request.Email}
		if err := db.Create(&participant).Error; err != nil {
			response := gin.H{"Message": "ERROR: BAD REQUEST"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		response := gin.H{"Message": "SUCCESS", "Data": participant}
		c.JSON(http.StatusCreated, response)
	}
}
