package models

import (
	"database/sql/driver"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type ParticipantCareerInterest string

const (
	SoftwareEngineering  ParticipantCareerInterest = "software-engineering"
	ProductManagement    ParticipantCareerInterest = "product-management"
	UIDesigner           ParticipantCareerInterest = "ui-designer"
	UXDesigner           ParticipantCareerInterest = "ux-designer"
	UXResearcher         ParticipantCareerInterest = "ux-researcher"
	ITConsultant         ParticipantCareerInterest = "it-consultant"
	GameDeveloper        ParticipantCareerInterest = "game-developer"
	CyberSecurity        ParticipantCareerInterest = "cyber-security"
	BusinessAnalyst      ParticipantCareerInterest = "business-analyst"
	BusinessIntelligence ParticipantCareerInterest = "business-intelligence"
	DataScientist        ParticipantCareerInterest = "data-scientist"
	DataAnalyst          ParticipantCareerInterest = "data-analyst"
)

func (participantCareerInterest *ParticipantCareerInterest) Scan(value interface{}) error {
	*participantCareerInterest = ParticipantCareerInterest(value.(string))
	return nil
}

func (participantCareerInterest ParticipantCareerInterest) Value() (driver.Value, error) {
	return string(participantCareerInterest), nil
}

type ParticipantStatus string

const (
	WaitingForVerification ParticipantStatus = "waiting-for-verification"
	Verified               ParticipantStatus = "verified"
	Declined               ParticipantStatus = "declined"
)

func (participantStatus *ParticipantStatus) Scan(value interface{}) error {
	*participantStatus = ParticipantStatus(value.(string))
	return nil
}

func (participantStatus ParticipantStatus) Value() (driver.Value, error) {
	return string(participantStatus), nil
}

type Participant struct {
	gorm.Model
	Name           string            `json:"name" gorm:"not null;unique"`
	Email          string            `json:"email" gorm:"not null;unique"`
	CareerInterest pq.StringArray    `json:"career_interest" gorm:"type:participant_career_interest[];default:array[]::participant_career_interest[];not null"`
	Status         ParticipantStatus `json:"status" gorm:"type:participant_status;not null"`
	Memberships    []Membership      `json:"memberships"`
	Photos         []Photo           `json:"photos"`
}
