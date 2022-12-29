package repository

import (
	"arkavidia-backend-8.0/competition/types"
)

type SignInTeamRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,ascii"`
}

type SignUpTeamRequest struct {
	Username string             `json:"username" binding:"required,alphanum"`
	Password string             `json:"password" binding:"required,ascii"`
	TeamName string             `json:"team_name" binding:"required,ascii"`
	Members  []SignUpMembership `json:"member_list" binding:"required,dive"`
}

type GetTeamQuery struct {
	TeamID uint `form:"team_id" field:"team_id" binding:"required,gt=0"`
}

type GetAllTeamsQuery struct {
	Page         int                `form:"page" field:"page" binding:"required,gt=0"`
	Size         int                `form:"size" field:"size" binding:"required,gt=0"`
	TeamCategory types.TeamCategory `form:"team_category" field:"team_category" binding:"required,oneof=competitive-programming datavidia uxvidia arkalogica"`
}

type ChangePasswordRequest struct {
	Password string `json:"password" binding:"required,ascii"`
}

type CompetitionRegistrationQuery struct {
	TeamCategory types.TeamCategory `form:"competition" field:"competition" binding:"required,oneof=competitive-programming datavidia uxvidia arkalogica"`
}

type ChangeStatusTeamQuery struct {
	TeamID uint `form:"team_id" field:"team_id" binding:"required,gt=0"`
}

type ChangeStatusTeamRequest struct {
	Status types.TeamStatus `json:"status" binding:"required,oneof=waiting-for-evaluation passed eliminated"`
}
