package models

import (
	"time"

	"github.com/google/uuid"
)

type SubmissionStage string

const (
	FirstStage  SubmissionStage = "first-stage"
	SecondStage SubmissionStage = "second-stage"
	FinalStage  SubmissionStage = "final-stage"
)

type Submission struct {
	ID            uuid.UUID       `json:"submission_id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TeamID        uuid.UUID       `json:"team_id" gorm:"type:uuid;not null"`
	FileName      uuid.UUID       `json:"link_to_file" gorm:"type:uuid;not null"`
	FileExtension string          `json:"file_extension" gorm:"not null"`
	Timestamp     time.Time       `json:"timestamp" gorm:"not null"`
	Stage         SubmissionStage `json:"stage" gorm:"not null;default:current_timestamp"`
	Team          Team            `json:"-" gorm:"foreignKey:TeamID;references:ID"`
}
