package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"

	"arkavidia-backend-8.0/competition/types"
)

type Participant struct {
	gorm.Model
	Name           string                           `gorm:"not null;unique"`
	Email          string                           `gorm:"not null;unique"`
	CareerInterest types.ParticipantCareerInterests `gorm:"not null"`
	Status         types.ParticipantStatus          `gorm:"not null"`
	Memberships    []Membership
	Photos         []Photo
}

type DisplayParticipant struct {
	ID             uint                             `json:"id,omitempty"`
	CreatedAt      time.Time                        `json:"created_at,omitempty"`
	UpdatedAt      time.Time                        `json:"updated_at,omitempty"`
	Name           string                           `json:"name,omitempty"`
	Email          string                           `json:"email,omitempty"`
	CareerInterest types.ParticipantCareerInterests `json:"career_interest,omitempty"`
	Status         types.ParticipantStatus          `json:"status,omitempty"`
	Memberships    []Membership                     `json:"memberships,omitempty"`
	Photos         []Photo                          `json:"photos,omitempty"`
}

func (participant Participant) MarshalJSON() ([]byte, error) {
	return json.Marshal(&DisplayParticipant{
		ID:             participant.ID,
		CreatedAt:      participant.CreatedAt,
		UpdatedAt:      participant.UpdatedAt,
		Name:           participant.Name,
		Email:          participant.Email,
		CareerInterest: participant.CareerInterest,
		Status:         participant.Status,
		Memberships:    participant.Memberships,
		Photos:         participant.Photos,
	})
}
