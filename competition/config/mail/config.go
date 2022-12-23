package mail

import (
	"os"
	"strconv"
)

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SenderName   string
	AuthEmail    string
	AuthPassword string
}

var currentEmailConfig *EmailConfig = nil

func Init() *EmailConfig {
	smtpHost := os.Getenv("CONFIG_SMTP_HOST")
	smtpPort, err := strconv.Atoi(os.Getenv("CONFIG_SMTP_PORT"))
	if err != nil {
		panic(err)
	}
	senderName := os.Getenv("CONFIG_SENDER_NAME")
	authEmail := os.Getenv("CONFIG_AUTH_EMAIL")
	authPassword := os.Getenv("CONFIG_AUTH_PASSWORD")

	return &EmailConfig{
		SMTPHost:     smtpHost,
		SMTPPort:     smtpPort,
		SenderName:   senderName,
		AuthEmail:    authEmail,
		AuthPassword: authPassword,
	}
}

func GetEmailConfig() *EmailConfig {
	if currentEmailConfig == nil {
		currentEmailConfig = Init()
	}
	return currentEmailConfig
}
