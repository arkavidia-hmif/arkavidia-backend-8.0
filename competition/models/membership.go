package models

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MembershipRole string

const (
	Leader    MembershipRole = "Leader"
	MemberOne MembershipRole = "Member 1"
	MemberTwo MembershipRole = "Member 2"
)

type Membership struct {
	TeamID        uuid.UUID      `json:"team_id" gorm:"type:uuid;primaryKey;uniqueIndex:composite_index"`
	ParticipantID uuid.UUID      `json:"participant_id" gorm:"type:uuid;primaryKey;uniqueIndex:composite_index"`
	Role          MembershipRole `json:"role" gorm:"not null;uniqueIndex:composite_index"`
	Team          Team           `json:"-" gorm:"foreignKey:TeamID;references:ID"`
	Participant   Participant    `json:"-" gorm:"foreignKey:ParticipantID;references:ID"`
}

// Menambahkan constraint untuk mengecek apakah terdapat participant yang mengikuti dua team atau lebih dengan jenis lomba yang sama atau memiliki role leader
func (membership *Membership) BeforeSave(tx *gorm.DB) error {
	conditionMembership := Membership{ParticipantID: membership.ParticipantID}
	newMemberships := []Membership{}
	if err := tx.Where(&conditionMembership).Find(&newMemberships).Error; err != nil {
		return err
	}

	for _, membershipA := range newMemberships {
		for _, membershipB := range newMemberships {
			if membershipA.TeamID != membershipB.TeamID {
				if membershipA.Team.TeamCategory == membershipB.Team.TeamCategory {
					return fmt.Errorf("Error: Invalid Database Operation!")
				}
				if membershipA.Role == Leader && membershipB.Role == Leader {
					return fmt.Errorf("Error: Invalid Database Operation!")
				}
			}
		}
	}

	return nil
}
