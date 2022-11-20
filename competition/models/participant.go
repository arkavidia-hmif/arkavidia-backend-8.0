package models

import (
	"github.com/google/uuid"
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

type Participant struct {
	gorm.Model
	ParticipantID  uuid.UUID      `json:"participant_id" gorm:"type:uuid;default:gen_random_uuid();unique"`
	Name           string         `json:"name" gorm:"not null;unique"`
	Email          string         `json:"email" gorm:"not null;unique"`
	CareerInterest pq.StringArray `json:"career_interest" gorm:"type:text[];default:array[]::text[];not null"`
}
