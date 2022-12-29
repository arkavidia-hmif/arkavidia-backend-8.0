package types

import (
	"database/sql/driver"
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
