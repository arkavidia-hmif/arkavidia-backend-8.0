package mail

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
	once    sync.Once
}

// Private
func (mailBroker *MailBroker) lazyInit() {
	mailBroker.once.Do(func() {
		config := messageConfig.Config.GetMetadata()

		// Asynchronous Channel
		mailBroker.channel = make(chan MailParameters, config.BufferSize)
	})
}

func (mailBroker *MailBroker) sendMailToClient(ctx context.Context, mailParameters MailParameters) error {
	mailBroker.lazyInit()

	// Synchronous Channel
	errorBroker := make(chan error)

	go func() {
		// TODO: Tambahkan SMTP menggunakan lib gomail
		// REFERENCE: https://dasarpemrogramangolang.novalagung.com/C-send-email.html
		// ASSIGNED TO: @samuelswandi
		// STATUS: IN PROGRESS

		// TODO: Gunakan templating HTML static file sebagai body email
		// REFERENCE: https://dasarpemrogramangolang.novalagung.com/B-template-render-html.html
		// ASSIGNED TO: @samuelswandi

		config := mailConfig.Config.GetMetadata()

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

func (mailBroker *MailBroker) recoverMailToBroker(mailParameters MailParameters) {
	mailBroker.lazyInit()

	if r := recover(); r != nil {
		mailBroker.AddMailToBroker(mailParameters)
	}
}

func (mailBroker *MailBroker) waitMailToClient(mailParameters MailParameters) {
	mailBroker.lazyInit()

	config := messageConfig.Config.GetMetadata()
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	err := mailBroker.sendMailToClient(ctx, mailParameters)

	defer mailBroker.recoverMailToBroker(mailParameters)
	if err != nil {
		panic(err)
	}
}

func (mailBroker *MailBroker) mailRun() {
	defer mailBroker.wg.Done()

	mailBroker.lazyInit()
	for mailParameters := range mailBroker.channel {
		mailBroker.waitMailToClient(mailParameters)
	}
}

// Public
func (mailBroker *MailBroker) AddMailToBroker(mailParameters MailParameters) {
	mailBroker.lazyInit()
	mailBroker.channel <- mailParameters
}

func (mailBroker *MailBroker) RunMailWorker(numOfWorkers int) {
	mailBroker.lazyInit()
	mailBroker.wg.Add(numOfWorkers)
	for i := 0; i < numOfWorkers; i++ {
		go mailBroker.mailRun()
	}
	mailBroker.wg.Wait()
}

func (mailBroker *MailBroker) CloseWorker() {
	mailBroker.lazyInit()
	close(mailBroker.channel)
}

var Broker = &MailBroker{}
