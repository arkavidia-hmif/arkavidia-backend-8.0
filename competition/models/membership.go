package models

import (
	"fmt"

	"gorm.io/gorm"

	"arkavidia-backend-8.0/competition/types"
)

type Membership struct {
	gorm.Model
	TeamID        uint                 `json:"team_id" gorm:"uniqueIndex:membership_index"`
	ParticipantID uint                 `json:"participant_id" gorm:"uniqueIndex:membership_index"`
	Role          types.MembershipRole `json:"role" gorm:"not null"`
	Team          Team                 `json:"team" gorm:"foreignKey:TeamID;references:ID"`
	Participant   Participant          `json:"participant" gorm:"foreignKey:ParticipantID;references:ID"`
}

// Menambahkan constraint untuk mengecek apakah terdapat participant yang mengikuti dua team atau lebih
// dengan jenis lomba yang sama atau memiliki role leader lebih dari satu kali
func (membership *Membership) BeforeSave(tx *gorm.DB) error {
	if membership.TeamID != 0 && membership.ParticipantID != 0 {
		conditionParticipantID := Membership{ParticipantID: membership.ParticipantID}
		oldMembershipsParticipantID := []Membership{}
		if err := tx.Where(&conditionParticipantID).Find(&oldMembershipsParticipantID).Error; err != nil {
			return err
		}

		for _, oldMembership := range oldMembershipsParticipantID {
			if oldMembership.Team.TeamCategory != "" && oldMembership.Team.TeamCategory == membership.Team.TeamCategory {
				return fmt.Errorf("ERROR: CANNOT PARTICIPATE MORE THAN ONCE")
			}
			if oldMembership.Role == types.Leader && membership.Role == types.Leader {
				return fmt.Errorf("ERROR: INELIGIBLE LEADER")
			}
		}

		conditionTeamID := Membership{TeamID: membership.TeamID}
		oldMembershipsTeamID := []Membership{}
		if err := tx.Where(&conditionTeamID).Find(&oldMembershipsTeamID).Error; err != nil {
			return err
		}

		for _, oldMembership := range oldMembershipsTeamID {
			if oldMembership.Role == types.Leader && membership.Role == types.Leader {
				return fmt.Errorf("ERROR: INELIGIBLE LEADER")
			}
		}
	}

	return nil
}
