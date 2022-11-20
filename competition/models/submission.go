package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubmissionStage string

const (
	FirstStage  SubmissionStage = "first-stage"
	SecondStage SubmissionStage = "second-stage"
	FinalStage  SubmissionStage = "final-stage"
)

type Submission struct {
	gorm.Model
	FileName      uuid.UUID       `json:"file_name" gorm:"type:uuid;unique"`
	FileExtension string          `json:"file_extension" gorm:"not null"`
	TeamID        uuid.UUID       `json:"team_id" gorm:"type:uuid;not null;uniqueIndex:submission_index"`
	Stage         SubmissionStage `json:"stage" gorm:"not null;default:current_timestamp;uniqueIndex:submission_index"`
	Team          Team            `json:"-" gorm:"foreignKey:TeamID;references:ID"`
}
