package worker

import (
	"fmt"
	"time"

	"gopkg.in/gomail.v2"

	mailConfig "arkavidia-backend-8.0/competition/config/mail"
	messageConfig "arkavidia-backend-8.0/competition/config/message"
	"arkavidia-backend-8.0/competition/utils/broker"
)

func SendMailToClient(mailParameters broker.MailParameters) error {
	// TODO: Tambahkan SMTP menggunakan lib gomail
	// REFERENCE: https://dasarpemrogramangolang.novalagung.com/C-send-email.html
	// ASSIGNED TO: @rayhankinan dan @samuelswandi

	config := mailConfig.GetEmailConfig()

	// TODO: Gunakan templating HTML static file sebagai body email
	// REFERENCE: https://dasarpemrogramangolang.novalagung.com/B-template-render-html.html
	// ASSIGNED TO: @rayhankinan dan @samuelswandi

	const subjectHeader = "Test Mail"
	const emailBody = "Hello, <b>have a nice day</b>"

	mailer := gomail.NewMessage()

	mailer.SetHeader("From", fmt.Sprintf("%s <%s>", config.SenderName, config.AuthEmail))
	mailer.SetHeader("To", mailParameters.Email)
	mailer.SetHeader("Subject", subjectHeader)
	mailer.SetAddressHeader("Cc", config.AuthEmail, config.SenderName)
	mailer.SetBody("text/html", emailBody)

	dialer := gomail.NewDialer(
		config.SMTPHost,
		config.SMTPPort,
		config.AuthEmail,
		config.AuthPassword,
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
