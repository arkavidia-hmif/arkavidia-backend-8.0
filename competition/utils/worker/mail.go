package worker

import (
	"context"
	"fmt"
	"sync"

	"gopkg.in/gomail.v2"

	mailConfig "arkavidia-backend-8.0/competition/config/mail"
	messageConfig "arkavidia-backend-8.0/competition/config/message"
)

type MailParameters struct {
	Email string
}

type MailBroker struct {
	channel chan MailParameters
	wg      sync.WaitGroup
}

var mailBroker *MailBroker = nil

func Init() *MailBroker {
	config := messageConfig.GetMessageConfig()

	// Asynchronous Channel
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

func SendMailToClient(ctx context.Context, mailParameters MailParameters) error {
	// TODO: Tambahkan SMTP menggunakan lib gomail
	// REFERENCE: https://dasarpemrogramangolang.novalagung.com/C-send-email.html
	// ASSIGNED TO: @rayhankinan dan @samuelswandi

	// TODO: Gunakan templating HTML static file sebagai body email
	// REFERENCE: https://dasarpemrogramangolang.novalagung.com/B-template-render-html.html
	// ASSIGNED TO: @rayhankinan dan @samuelswandi

	// Synchronous Channel
	errorBroker := make(chan error)

	go func() {
		config := mailConfig.GetEmailConfig()

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

		errorBroker <- dialer.DialAndSend(mailer)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errorBroker:
		return err
	}
}

func RecoverMailToBroker(mailParameters MailParameters) {
	if r := recover(); r != nil {
		AddMailToBroker(mailParameters)
	}
}

func WaitMailToClient(mailParameters MailParameters) {
	config := messageConfig.GetMessageConfig()

	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	defer RecoverMailToBroker(mailParameters)

	err := SendMailToClient(ctx, mailParameters)
	if err != nil {
		panic(err)
	}
}

func AddMailToBroker(mailParameters MailParameters) {
	mailBroker := GetMailBroker()

	mailBroker.channel <- mailParameters
}

func MailRun() {
	mailBroker := GetMailBroker()

	defer mailBroker.wg.Done()
	for mailParameters := range mailBroker.channel {
		WaitMailToClient(mailParameters)
	}
}

func RunMailWorker(numOfWorkers int) {
	mailBroker := GetMailBroker()

	mailBroker.wg.Add(numOfWorkers)
	for i := 0; i < numOfWorkers; i++ {
		go MailRun()
	}
	mailBroker.wg.Wait()
}

func CloseWorker() {
	mailBroker := GetMailBroker()

	close(mailBroker.channel)
}
