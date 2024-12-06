package mail

import (
	"net/smtp"
)

type SmtpEmailSender struct {
	fromEmailAddress  string
	fromEmailPassword string
}

func SmtpNewGmailSender(fromEmailAddress, fromEmailPassword string) *SmtpEmailSender {
	return &SmtpEmailSender{
		fromEmailAddress:  fromEmailAddress,
		fromEmailPassword: fromEmailPassword,
	}
}

func (emailSender *SmtpEmailSender) sendEmailSmtp(to []string, content string) error {

	// Message.
	message := []byte(content)

	// Authentication.
	auth := smtp.PlainAuth("", emailSender.fromEmailAddress, emailSender.fromEmailPassword, smtpAuthAddress)

	// Sending email.
	err := smtp.SendMail(smtpServerAddress, auth, emailSender.fromEmailAddress, to, message)

	return err
}
