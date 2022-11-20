package models

import (
	"database/sql/driver"

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
	*participantCareerInterest = ParticipantCareerInterest(value.([]byte))
	return nil
}

func (participantCareerInterest ParticipantCareerInterest) Value() (driver.Value, error) {
	return string(participantCareerInterest), nil
}

type Participant struct {
	gorm.Model
	Name           string                      `json:"name" gorm:"not null;unique"`
	Email          string                      `json:"email" gorm:"not null;unique"`
	CareerInterest []ParticipantCareerInterest `json:"career_interest" gorm:"type:participant_career_interest[];default:array[]::participant_career_interest[];not null"`
	Memberships    []Membership                `json:"memberships"`
	Photos         []Photo                     `json:"photos"`
}
