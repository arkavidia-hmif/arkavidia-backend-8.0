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
	*participantCareerInterest = ParticipantCareerInterest(value.(string))
	return nil
}

func (participantCareerInterest ParticipantCareerInterest) Value() (driver.Value, error) {
	return string(participantCareerInterest), nil
}

func (ParticipantCareerInterest) GormDataType() string {
	return "participant_career_interest"
}

type ParticipantCareerInterests []ParticipantCareerInterest

func (participantCareerInterests *ParticipantCareerInterests) Scan(values []interface{}) error {
	for _, value := range values {
		*participantCareerInterests = append(*participantCareerInterests, ParticipantCareerInterest(value.(string)))
	}
	return nil
}

func (participantCareerInterests ParticipantCareerInterests) Value() (driver.Value, error) {
	var values []string
	for _, participantCareerInterest := range participantCareerInterests {
		values = append(values, string(participantCareerInterest))
	}
	return values, nil
}

func (ParticipantCareerInterests) GormDataType() string {
	return "participant_career_interest[]"
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

func (ParticipantStatus) GormDataType() string {
	return "participant_status"
}

type Participant struct {
	gorm.Model
	Name           string                     `json:"name" gorm:"not null;unique"`
	Email          string                     `json:"email" gorm:"not null;unique"`
	CareerInterest ParticipantCareerInterests `json:"career_interest" gorm:"not null"`
	Status         ParticipantStatus          `json:"status" gorm:"not null"`
	Memberships    []Membership               `json:"memberships"`
	Photos         []Photo                    `json:"photos"`
}
