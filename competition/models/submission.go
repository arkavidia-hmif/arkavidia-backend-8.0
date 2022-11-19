package models

import (
	"time"

	"github.com/google/uuid"
)

type SubmissionStage string

const (
	FirstStage  SubmissionStage = "First Stage"
	SecondStage SubmissionStage = "Second Stage"
	FinalStage  SubmissionStage = "Final Stage"
)

type Submission struct {
	ID         uuid.UUID       `json:"submission_id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TeamID     uuid.UUID       `json:"team_id" gorm:"type:uuid;not null"`
	LinkToFile uuid.UUID       `json:"link_to_file" gorm:"type:uuid;not null"`
	Timestamp  time.Time       `json:"timestamp" gorm:"not null"`
	Stage      SubmissionStage `json:"stage" gorm:"not null;default:current_timestamp"`
	Team       Team            `json:"-" gorm:"foreignKey:TeamID;references:ID"`
}
