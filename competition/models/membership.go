package models

import (
	"database/sql/driver"
	"fmt"

	"gorm.io/gorm"
)

type MembershipRole string

const (
	Leader MembershipRole = "leader"
	Member MembershipRole = "member"
)

func (membershipRole *MembershipRole) Scan(value interface{}) error {
	*membershipRole = MembershipRole(value.(string))
	return nil
}

func (membershipRole MembershipRole) Value() (driver.Value, error) {
	return string(membershipRole), nil
}

func (MembershipRole) GormDataType() string {
	return "membership_role"
}

type Membership struct {
	gorm.Model
	TeamID        uint           `json:"team_id" gorm:"uniqueIndex:membership_index"`
	ParticipantID uint           `json:"participant_id" gorm:"uniqueIndex:membership_index"`
	Role          MembershipRole `json:"role" gorm:"not null"`
	Team          Team           `json:"team" gorm:"foreignKey:TeamID;references:ID"`
	Participant   Participant    `json:"participant" gorm:"foreignKey:ParticipantID;references:ID"`
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
			if oldMembership.Role == Leader && membership.Role == Leader {
				return fmt.Errorf("ERROR: INELIGIBLE LEADER")
			}
		}

		conditionTeamID := Membership{TeamID: membership.TeamID}
		oldMembershipsTeamID := []Membership{}
		if err := tx.Where(&conditionTeamID).Find(&oldMembershipsTeamID).Error; err != nil {
			return err
		}

		for _, oldMembership := range oldMembershipsTeamID {
			if oldMembership.Role == Leader && membership.Role == Leader {
				return fmt.Errorf("ERROR: INELIGIBLE LEADER")
			}
		}
	}

	return nil
}
