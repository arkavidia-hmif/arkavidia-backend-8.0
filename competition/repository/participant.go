package repository

import (
	"arkavidia-backend-8.0/competition/types"
)

type GetMemberQuery struct {
	TeamID uint `form:"team_id" field:"team_id" binding:"required,gt=0"`
}

type GetAllMembersQuery struct {
	Page int `form:"page" field:"page" binding:"required,gt=0"`
	Size int `form:"size" field:"size" binding:"required,gt=0"`
}

type AddMemberRequest struct {
	Name            string                           `json:"name" binding:"required,ascii"`
	Email           string                           `json:"email" binding:"required,email"`
	CareerInterests types.ParticipantCareerInterests `json:"career_interest" binding:"required,dive,oneof=software-engineering product-management ui-designer ux-designer ux-researcher it-consultant game-developer cyber-security business-analyst business-intelligence data-scientist data-analyst"`
	Role            types.MembershipRole             `json:"role" binding:"required,oneof=leader member"`
}

type ChangeCareerInterestQuery struct {
	ParticipantID uint `form:"participant_id" field:"participant_id" binding:"required,gt=0"`
}

type ChangeCareerInterestRequest struct {
	CareerInterests types.ParticipantCareerInterests `json:"career_interest" binding:"required,dive,oneof=software-engineering product-management ui-designer ux-designer ux-researcher it-consultant game-developer cyber-security business-analyst business-intelligence data-scientist data-analyst"`
}

type ChangeRoleQuery struct {
	ParticipantID uint `form:"participant_id" field:"participant_id" binding:"required,gt=0"`
}

type ChangeRoleRequest struct {
	Role types.MembershipRole `json:"role" binding:"required,oneof=leader member"`
}

type ChangeStatusParticipantQuery struct {
	ParticipantID uint `form:"participant_id" field:"participant_id" binding:"required,gt=0"`
}

type ChangeStatusParticipantRequest struct {
	Status types.ParticipantStatus `json:"status" binding:"required,oneof=waiting-for-verification verified declined"`
}

type DeleteMemberRequest struct {
	ParticipantID uint `json:"participant_id" binding:"required,gt=0"`
}
