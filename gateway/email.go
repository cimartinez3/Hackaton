package gateway

import (
	"bytes"
	"fmt"
	"net/smtp"
	"text/template"
)

type IEmail interface {
	SendOTP(email string) error
}

type Email struct {
}

func NewEmailService() IEmail {
	return &Email{}
}

func (e *Email) SendOTP(email string) error {
	from := "theboyshackaton@gmail.com"
	password := "smnlojgcmhdzelqo"

	// Receiver email address.
	to := []string{
		email,
	}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	t, err := template.ParseFiles("./gateway/template.html")

	if err != nil {
		fmt.Println(err)
		return err
	}

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: OTP Verification \n%s\n\n", mimeHeaders)))

	t.Execute(&body, struct{}{})

	// Sending email.
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("Email Sent!")
	return nil
}
