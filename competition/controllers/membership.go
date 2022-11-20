package controllers

import (
	"net/http"

	"arkavidia-backend-8.0/competition/models"
	databaseService "arkavidia-backend-8.0/competition/services/database"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MembershipRequest struct {
	TeamID        uuid.UUID             `json:"team_id" gorm:"type:uuid;uniqueIndex:membership_index"`
	ParticipantID uuid.UUID             `json:"participant_id" gorm:"type:uuid;uniqueIndex:membership_index"`
	Role          models.MembershipRole `json:"role" gorm:"type:membership_role;not null"`
	Team          models.Team           `json:"team" gorm:"foreignKey:TeamID;references:UUID"`
	Participant   models.Participant    `json:"participant" gorm:"foreignKey:ParticipantID;references:UUID"`
}

func GetMember() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.GetDB()
		teamID := c.MustGet("team_id").(uuid.UUID)

		membership := models.Membership{TeamID: teamID}
		if err := db.Preload("Memberships").Find(&membership).Error; err != nil {
			response := gin.H{"Message": "ERROR: BAD REQUEST"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		response := gin.H{"Message": "SUCCESS", "Data": membership}
		c.JSON(http.StatusOK, response)

	}
}

func AddMembership() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := databaseService.GetDB()
		request := MembershipRequest{}

		if err := c.BindJSON(&request); err != nil {
			response := gin.H{"Message": "ERROR: BAD REQUEST"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		membership := models.Membership{
			TeamID:        request.TeamID,
			ParticipantID: request.ParticipantID,
			Role:          request.Role,
			Team:          request.Team,
			Participant:   request.Participant,
		}
		if err := db.Create(&membership).Error; err != nil {
			response := gin.H{"Message": "ERROR: BAD REQUEST"}
			c.JSON(http.StatusBadRequest, response)
			return
		}

		response := gin.H{"Message": "SUCCESS"}
		c.JSON(http.StatusCreated, response)
	}
}
