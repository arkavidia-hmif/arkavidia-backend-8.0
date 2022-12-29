package repository

import (
	"mime/multipart"

	"arkavidia-backend-8.0/competition/types"
)

type GetPhotoQuery struct {
	ParticipantID uint `form:"participant_id" field:"participant_id" binding:"required,gt=0"`
}

type GetAllPhotosQuery struct {
	Page int `form:"page" field:"page" binding:"required,gt=0"`
	Size int `form:"size" field:"size" binding:"required,gt=0"`
}

type AdminDownloadPhotoQuery struct {
	PhotoID uint `form:"photo_id" field:"photo_id" binding:"required,gt=0"`
}

type TeamDownloadPhotoQuery struct {
	ParticipantID uint `form:"participant_id" field:"participant_id" binding:"required,gt=0"`
	PhotoID       uint `form:"photo_id" field:"photo_id" binding:"required,gt=0"`
}

type AddPhotoRequest struct {
	ParticipantID uint                  `form:"participant_id" field:"participant_id" binding:"required,gt=0"`
	Type          types.PhotoType       `form:"type" field:"type" binding:"required,oneof=pribadi kartu-pelajar bukti-mahasiswa-aktif bukti-pembayaran"`
	File          *multipart.FileHeader `form:"file" field:"file" binding:"required"`
}

type ChangeStatusPhotoQuery struct {
	PhotoID uint `form:"photo_id" field:"photo_id" binding:"required,gt=0"`
}

type ChangeStatusPhotoRequest struct {
	Status types.PhotoStatus `json:"status" binding:"required,oneof=waiting-for-approval approved denied"`
}

type DeletePhotoRequest struct {
	ParticipantID uint   `json:"participant_id" binding:"required,gt=0"`
	FileName      string `json:"file_name" binding:"required,uuid"`
	FileExtension string `json:"file_extension" binding:"required,alpha"`
}
