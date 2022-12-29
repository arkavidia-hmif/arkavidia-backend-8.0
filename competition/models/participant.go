package models

import (
	"gorm.io/gorm"

	"arkavidia-backend-8.0/competition/types"
)

type Participant struct {
	gorm.Model
	Name           string                           `json:"name" gorm:"not null;unique"`
	Email          string                           `json:"email" gorm:"not null;unique"`
	CareerInterest types.ParticipantCareerInterests `json:"career_interest" gorm:"not null"`
	Status         types.ParticipantStatus          `json:"status" gorm:"not null"`
	Memberships    []Membership                     `json:"memberships"`
	Photos         []Photo                          `json:"photos"`
}
