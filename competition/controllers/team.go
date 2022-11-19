package controllers

import (
	"github.com/gin-gonic/gin"
)

type SignInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignUpRequest struct {
	Username          string   `json:"username"`
	Password          string   `json:"password"`
	ListOfMemberEmail []string `json:"email_list"`
}

func SignInHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
