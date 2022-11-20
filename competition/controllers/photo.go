package controllers

import (
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"arkavidia-backend-8.0/competition/models"
)

type GetPhotoRequest struct {
	ParticipantID uuid.UUID `form:"participant_id" field:"participant_id"`
}

type AddPhotoRequest struct {
	ParticipantID uuid.UUID             `form:"participant_id" field:"participant_id"`
	Type          models.PhotoType      `form:"type" field:"type"`
	File          *multipart.FileHeader `form:"file" field:"file" binding:"required"`
}

type DeletePhotoRequest struct {
	ParticipantID uuid.UUID             `form:"participant_id" field:"participant_id"`
	Type          models.PhotoType      `form:"type" field:"type"`
	File          *multipart.FileHeader `form:"file" field:"file" binding:"required"`
}

func GetPhotoHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
