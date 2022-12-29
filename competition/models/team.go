package models

import (
	"fmt"

	"gorm.io/gorm"

	"arkavidia-backend-8.0/competition/types"
)

type Team struct {
	gorm.Model
	Username       string                `json:"username" gorm:"not null;unique"`
	HashedPassword types.EncryptedString `json:"password" gorm:"not null" visibility:"false"`
	TeamName       string                `json:"team_name" gorm:"not null;unique"`
	TeamCategory   types.TeamCategory    `json:"team_category" gorm:"default:null"`
	AdminID        uint                  `json:"admin_id" gorm:"default:null"`
	Status         types.TeamStatus      `json:"status" gorm:"not null"`
	Memberships    []Membership          `json:"memberships"`
	Submissions    []Submission          `json:"submissions"`
	ApprovedBy     Admin                 `json:"admin" gorm:"foreignKey:AdminID;references:ID"`
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
