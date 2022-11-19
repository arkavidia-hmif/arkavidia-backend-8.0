package models

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type ParticipantCareerInterest string

const (
	SoftwareEngineering  ParticipantCareerInterest = "Software Engineering"
	ProductManagement    ParticipantCareerInterest = "Product Management"
	UIDesigner           ParticipantCareerInterest = "UI Designer"
	UXDesigner           ParticipantCareerInterest = "UX Designer"
	UXResearcher         ParticipantCareerInterest = "UX Researcher"
	ITConsultant         ParticipantCareerInterest = "IT Consultant"
	GameDeveloper        ParticipantCareerInterest = "Game Developer"
	CyberSecurity        ParticipantCareerInterest = "Cyber Security"
	BusinessAnalyst      ParticipantCareerInterest = "Business Analyst"
	BusinessIntelligence ParticipantCareerInterest = "Business Intelligence"
	DataScientist        ParticipantCareerInterest = "Data Scientist"
	DataAnalyst          ParticipantCareerInterest = "Data Analyst"
)

type Participant struct {
	ID             uuid.UUID      `json:"participant_id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name           string         `json:"name" gorm:"not null;unique"`
	Email          string         `json:"email" gorm:"not null;unique"`
	CareerInterest pq.StringArray `json:"career_interest" gorm:"type:text[];not null"`
}
