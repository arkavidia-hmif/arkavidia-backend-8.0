package controllers

import (
	"arkavidia-backend-8.0/competition/models"
)

type SignUpMembership struct {
	Name            string                            `json:"name" binding:"required,ascii"`
	Email           string                            `json:"email" binding:"required,email"`
	CareerInterests models.ParticipantCareerInterests `json:"career_interest" binding:"required,dive,oneof=software-engineering product-management ui-designer ux-designer ux-researcher it-consultant game-developer cyber-security business-analyst business-intelligence data-scientist data-analyst"`
	Role            models.MembershipRole             `json:"role" binding:"required,oneof=leader member"`
}
