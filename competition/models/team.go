package models

import (
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"

	"arkavidia-backend-8.0/competition/types"
)

type Team struct {
	gorm.Model
	Username       string                `gorm:"not null;unique"`
	HashedPassword types.EncryptedString `gorm:"not null"`
	TeamName       string                `gorm:"not null;unique"`
	TeamCategory   types.TeamCategory    `gorm:"default:null"`
	AdminID        uint                  `gorm:"default:null"`
	Status         types.TeamStatus      `gorm:"not null"`
	ApprovedBy     Admin                 `gorm:"foreignKey:AdminID;references:ID"`
	Memberships    []Membership
	Submissions    []Submission
}

type DisplayTeam struct {
	ID           uint               `json:"id,omitempty"`
	CreatedAt    time.Time          `json:"created_at,omitempty"`
	UpdatedAt    time.Time          `json:"updated_at,omitempty"`
	Username     string             `json:"username,omitempty"`
	TeamName     string             `json:"team_name,omitempty"`
	TeamCategory types.TeamCategory `json:"team_category,omitempty"`
	AdminID      uint               `json:"admin_id,omitempty"`
	Status       types.TeamStatus   `json:"status,omitempty"`
	Memberships  []Membership       `json:"memberships,omitempty"`
	Submissions  []Submission       `json:"submissions,omitempty"`
}

func (team Team) MarshalJSON() ([]byte, error) {
	return json.Marshal(&DisplayTeam{
		ID:           team.ID,
		CreatedAt:    team.CreatedAt,
		UpdatedAt:    team.UpdatedAt,
		Username:     team.Username,
		TeamName:     team.TeamName,
		TeamCategory: team.TeamCategory,
		AdminID:      team.AdminID,
		Status:       team.Status,
		Memberships:  team.Memberships,
		Submissions:  team.Submissions,
	})
}

// Menambahkan constraint untuk mengecek apakah terdapat photos yang telah diapprove namun admin tidak tercatat
// atau photos yang belum diapprove namun admin tercatat
func (team *Team) BeforeSave(tx *gorm.DB) error {
	if team.Status != "" {
		if team.Status != types.WaitingForEvaluation && team.AdminID == 0 {
			return fmt.Errorf("ERROR: ADMIN MUST BE RECORDED")
		}
		if team.Status == types.WaitingForEvaluation && team.AdminID != 0 {
			return fmt.Errorf("ERROR: STATUS MUST BE RECORDED")
		}
	}

	return nil
}
