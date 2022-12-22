package worker

import (
	"time"

	messageConfig "arkavidia-backend-8.0/competition/config/message"
	"arkavidia-backend-8.0/competition/utils/broker"
)

func SendMailToClient(mailParameters broker.MailParameters) error {
	// TODO: Tambahkan SMTP menggunakan lib gomail
	// REFERENCE: https://dasarpemrogramangolang.novalagung.com/C-send-email.html
	// ASSIGNED TO: @samuelswandi

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
