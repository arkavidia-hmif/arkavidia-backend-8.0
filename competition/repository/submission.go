package repository

import (
	"mime/multipart"

	"arkavidia-backend-8.0/competition/types"
)

type GetSubmissionQuery struct {
	TeamID uint `form:"team_id" field:"team_id" binding:"required,gt=0"`
}

type GetAllSubmissionsQuery struct {
	Page int `form:"page" field:"page" binding:"required,gt=0"`
	Size int `form:"size" field:"size" binding:"required,gt=0"`
}

type DownloadSubmissionQuery struct {
	SubmissionID uint `form:"submission_id" field:"submission_id" binding:"required,gt=0"`
}

type AddSubmissionRequest struct {
	Stage types.SubmissionStage `form:"stage" field:"stage" binding:"required,oneof=first-stage second-stage final-stage"`
	File  *multipart.FileHeader `form:"file" field:"file" binding:"required"`
}

type DeleteSubmissionRequest struct {
	FileName      string `json:"file_name" binding:"required,uuid"`
	FileExtension string `json:"file_extension" binding:"required,alpha"`
}
