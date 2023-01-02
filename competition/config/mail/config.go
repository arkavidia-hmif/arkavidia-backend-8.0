package mail

import (
	"os"
	"strconv"
	"sync"
)

type EmailMetadata struct {
	SMTPHost     string
	SMTPPort     int
	SenderName   string
	AuthEmail    string
	AuthPassword string
}

type EmailConfig struct {
	metadata EmailMetadata
	once     sync.Once
}

// Private
func (emailConfig *EmailConfig) lazyInit() {
	emailConfig.once.Do(func() {
		smtpHost := os.Getenv("CONFIG_SMTP_HOST")
		smtpPort, err := strconv.Atoi(os.Getenv("CONFIG_SMTP_PORT"))
		if err != nil {
			panic(err)
		}
		senderName := os.Getenv("CONFIG_SENDER_NAME")
		authEmail := os.Getenv("CONFIG_AUTH_EMAIL")
		authPassword := os.Getenv("CONFIG_AUTH_PASSWORD")

		emailConfig.metadata.SMTPHost = smtpHost
		emailConfig.metadata.SMTPPort = smtpPort
		emailConfig.metadata.SenderName = senderName
		emailConfig.metadata.AuthEmail = authEmail
		emailConfig.metadata.AuthPassword = authPassword
	})
}

// Public
func (emailConfig *EmailConfig) GetMetadata() EmailMetadata {
	emailConfig.lazyInit()
	return emailConfig.metadata
}

var Config = &EmailConfig{}
