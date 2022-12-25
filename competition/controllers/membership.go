package controllers

import (
	"github.com/lib/pq"

	"arkavidia-backend-8.0/competition/models"
)

type SignUpMembership struct {
	Name           string                `json:"name" binding:"required"`
	Email          string                `json:"email" binding:"required"`
	CareerInterest pq.StringArray        `json:"career_interest" binding:"required"`
	Role           models.MembershipRole `json:"role" binding:"required"`
}
