package worker

import (
	"time"

	"gopkg.in/gomail.v2"

	messageConfig "arkavidia-backend-8.0/competition/config/message"
	"arkavidia-backend-8.0/competition/utils/broker"
)

const CONFIG_SMTP_HOST = "smtp.gmail.com"
const CONFIG_SMTP_PORT = 587
const CONFIG_SENDER_NAME = "Arkavidia <admin_arkavidia@gmail.com>"
const CONFIG_AUTH_EMAIL = "emailanda@gmail.com"
const CONFIG_AUTH_PASSWORD = "passwordemailanda"

func SendMailToClient(mailParameters broker.MailParameters) error {
	// TODO: Tambahkan SMTP menggunakan lib gomail
	// REFERENCE: https://dasarpemrogramangolang.novalagung.com/C-send-email.html
	// ASSIGNED TO: @samuelswandi

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", CONFIG_SENDER_NAME)
	mailer.SetHeader("To", mailParameters.Email)
	mailer.SetAddressHeader("Cc", "admin_arkavidia@gmail.com", "Admin Arkavidia")
	mailer.SetHeader("Subject", "Test mail")
	mailer.SetBody("text/html", "Hello, <b>have a nice day</b>")

	dialer := gomail.NewDialer(
		CONFIG_SMTP_HOST,
		CONFIG_SMTP_PORT,
		CONFIG_AUTH_EMAIL,
		CONFIG_AUTH_PASSWORD,
	)

	err := dialer.DialAndSend(mailer)

	return err
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
