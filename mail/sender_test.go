package mail

import (
	"testing"

	"github.com/Yelsnik/e-commerce-api/util"
	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T) {
	config, err := util.LoadConfig("..")

	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "A test email"
	content := `
	<h1>Hello World </h1>
	<p> This is a test meassage</p>
	`

	to := []string{"kingsleyokgeorge@gmail.com"}
	attachFiles := []string{}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	//require.NoError(t, err)
}

func TestSendEmailWithSMTP(t *testing.T) {
	config, err := util.LoadConfig("..")

	require.NoError(t, err)

	senderSmtp := SmtpNewGmailSender(config.EmailSenderAddress, config.EmailSenderPassword)

	content := "A test email \n this is a test email"
	to := []string{"kingsleyokgeorge@gmail.com"}

	err = senderSmtp.sendEmailSmtp(to, content)
	//require.NoError(t, err)
}
