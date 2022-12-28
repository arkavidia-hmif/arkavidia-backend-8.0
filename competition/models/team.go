package models

import (
	"database/sql/driver"
	"fmt"

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
	*teamCategory = TeamCategory(value.(string))
	return nil
}

func (teamCategory TeamCategory) Value() (driver.Value, error) {
	return string(teamCategory), nil
}

type TeamStatus string

const (
	WaitingForEvaluation TeamStatus = "waiting-for-evaluation"
	Passed               TeamStatus = "passed"
	Eliminated           TeamStatus = "eliminated"
)

func (teamStatus *TeamStatus) Scan(value interface{}) error {
	*teamStatus = TeamStatus(value.(string))
	return nil
}

func (teamStatus TeamStatus) Value() (driver.Value, error) {
	return string(teamStatus), nil
}

type Team struct {
	gorm.Model
	Username       string       `json:"username" gorm:"not null;unique"`
	HashedPassword []byte       `json:"password" gorm:"not null" visibility:"false"`
	TeamName       string       `json:"team_name" gorm:"not null;unique"`
	TeamCategory   TeamCategory `json:"team_category" gorm:"type:team_category;default:null"`
	AdminID        uint         `json:"admin_id" gorm:"default:null"`
	Status         TeamStatus   `json:"status" gorm:"type:team_status;not null"`
	Memberships    []Membership `json:"memberships"`
	Submissions    []Submission `json:"submissions"`
	ApprovedBy     Admin        `json:"admin" gorm:"foreignKey:AdminID;references:ID"`
}

// Menambahkan constraint untuk mengecek apakah terdapat photos yang telah diapprove namun admin tidak tercatat
// atau photos yang belum diapprove namun admin tercatat
func (team *Team) BeforeSave(tx *gorm.DB) error {
	if team.Status != "" {
		if team.Status != WaitingForEvaluation && team.AdminID == 0 {
			return fmt.Errorf("ERROR: ADMIN MUST BE RECORDED")
		}
		if team.Status == WaitingForEvaluation && team.AdminID != 0 {
			return fmt.Errorf("ERROR: STATUS MUST BE RECORDED")
		}
	}

	return nil
}
