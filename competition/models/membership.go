package models

import (
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"

	"arkavidia-backend-8.0/competition/types"
)

type Membership struct {
	gorm.Model
	TeamID        uint                 `gorm:"uniqueIndex:membership_index"`
	ParticipantID uint                 `gorm:"uniqueIndex:membership_index"`
	Role          types.MembershipRole `gorm:"not null"`
	Team          Team                 `gorm:"foreignKey:TeamID;references:ID"`
	Participant   Participant          `gorm:"foreignKey:ParticipantID;references:ID"`
}

type DisplayMembership struct {
	ID            uint                 `json:"id,omitempty"`
	CreatedAt     time.Time            `json:"created_at,omitempty"`
	UpdatedAt     time.Time            `json:"updated_at,omitempty"`
	TeamID        uint                 `json:"team_id,omitempty"`
	ParticipantID uint                 `json:"participant_id,omitempty"`
	Role          types.MembershipRole `json:"role,omitempty"`
}

func (membership Membership) MarshalJSON() ([]byte, error) {
	return json.Marshal(&DisplayMembership{
		ID:            membership.ID,
		CreatedAt:     membership.CreatedAt,
		UpdatedAt:     membership.UpdatedAt,
		TeamID:        membership.TeamID,
		ParticipantID: membership.ParticipantID,
		Role:          membership.Role,
	})
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
