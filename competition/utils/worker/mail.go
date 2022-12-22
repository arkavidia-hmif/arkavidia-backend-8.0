package worker

import (
	"time"

	messageConfig "arkavidia-backend-8.0/competition/config/message"
	"arkavidia-backend-8.0/competition/utils/broker"
)

func SendMailToClient(mailParameters broker.MailParameters) error {
	// TODO: Tambahkan SMTP pake lib gomail di sini @StaffBE
	// REFERENCE: https://dasarpemrogramangolang.novalagung.com/C-send-email.html

	return nil
}

func MailRun() {
	config := messageConfig.GetMessageConfig()
	mailBroker := broker.GetMailBroker()

	for {
		if broker.GetLength(mailBroker) > 0 {
			mailParameters := broker.ReceiveMailTask(mailBroker)
			go SendMailToClient(mailParameters)
		} else {
			time.Sleep(config.ReloadTime)
		}
	}
}
