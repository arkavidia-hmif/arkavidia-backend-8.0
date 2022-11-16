package models

import (
	"github.com/google/uuid"
)

type MembershipRole string

const (
	Leader    MembershipRole = "Leader"
	MemberOne MembershipRole = "Member 1"
	MemberTwo MembershipRole = "Member 2"
)

type Membership struct {
	TeamID        uuid.UUID      `json:"team_id" gorm:"type:uuid;primaryKey"`
	ParticipantID uuid.UUID      `json:"participant_id" gorm:"type:uuid;primaryKey"`
	Role          MembershipRole `json:"role" gorm:"not null"`
	Team          Team           `json:"-" gorm:"foreignKey:TeamID;references:ID"`
	Participant   Participant    `json:"-" gorm:"foreignKey:ParticipantID;references:ID"`
}
