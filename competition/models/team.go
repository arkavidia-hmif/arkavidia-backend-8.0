package models

import (
	"github.com/google/uuid"
)

type TeamCategory string

const (
	CP         TeamCategory = "cp"
	Datavidia  TeamCategory = "datavidia"
	UXVidia    TeamCategory = "uxvidia"
	Arkalogica TeamCategory = "arkalogica"
)

type Team struct {
	ID             uuid.UUID    `json:"team_id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Username       string       `json:"username" gorm:"not null;unique"`
	HashedPassword []byte       `json:"password" gorm:"not null"`
	TeamName       string       `json:"team_name" gorm:"not null;unique"`
	TeamCategory   TeamCategory `json:"team_category"`
}
