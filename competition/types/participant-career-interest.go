package types

import (
	"database/sql/driver"
	"regexp"
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

func (participantCareerInterests *ParticipantCareerInterests) Scan(values interface{}) error {
	regex, err := regexp.Compile(`[a-zA-Z\-]+`)
	if err != nil {
		return nil
	}

	words := regex.FindAllString(values.(string), -1)
	*participantCareerInterests = []ParticipantCareerInterest{}
	for _, word := range words {
		*participantCareerInterests = append(*participantCareerInterests, ParticipantCareerInterest(word))
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
