package models

import (
	"database/sql/driver"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TeamCategory string

const (
	CP         TeamCategory = "competitive-programming"
	Datavidia  TeamCategory = "datavidia"
	UXVidia    TeamCategory = "uxvidia"
	Arkalogica TeamCategory = "arkalogica"
)

func (teamCategory *TeamCategory) Scan(value interface{}) error {
	*teamCategory = TeamCategory(value.([]byte))
	return nil
}

func (teamCategory TeamCategory) Value() (driver.Value, error) {
	return string(teamCategory), nil
}

type Team struct {
	gorm.Model
	UUID           uuid.UUID    `json:"team_id" gorm:"type:uuid;not null;unique"`
	Username       string       `json:"username" gorm:"not null;unique"`
	HashedPassword []byte       `json:"password" gorm:"not null"`
	TeamName       string       `json:"team_name" gorm:"not null;unique"`
	TeamCategory   TeamCategory `json:"team_category" gorm:"type:team_category"`
}

func (team *Team) BeforeCreate(tx *gorm.DB) error {
	team.UUID = uuid.New()
	return nil
}
