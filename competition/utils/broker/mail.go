package broker

import (
	messageConfig "arkavidia-backend-8.0/competition/config/message"
)

type MailParameters struct {
	Email string
}

type MailBroker struct {
	channel chan MailParameters
}

var mailBroker *MailBroker = nil

func Init() *MailBroker {
	config := messageConfig.GetMessageConfig()

	return &MailBroker{
		channel: make(chan MailParameters, config.BufferSize),
	}
}

func GetMailBroker() *MailBroker {
	if mailBroker == nil {
		mailBroker = Init()
	}

	return mailBroker
}

func GetLength(mailBroker *MailBroker) int {
	return len(mailBroker.channel)
}

func SendMailTask(mailBroker *MailBroker, mailParameters MailParameters) {
	mailBroker.channel <- mailParameters
}

func ReceiveMailTask(mailBroker *MailBroker) MailParameters {
	mailParameters := <-mailBroker.channel

	return mailParameters
}
