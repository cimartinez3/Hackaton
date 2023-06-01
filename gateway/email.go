package gateway

import (
	"bytes"
	"fmt"
	"net/smtp"
	"text/template"
)

type IEamil interface {
	SendOTP(email string) bool
}

type Email struct {
}

func NewEmailService() IEamil {
	return &Email{}
}

func (e *Email) SendOTP(email string) bool {
	// Sender data.
	from := "theboyshackaton@hotmail.com"
	password := "_96notakcahsyobeht"

	// Receiver email address.
	to := []string{
		"carlos.martinez@kushkipagos.com",
	}

	// smtp server configuration.
	smtpHost := "smtp.office365.com"
	smtpPort := "587"

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	t, _ := template.ParseFiles("template.html")

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: This is a test subject \n%s\n\n", mimeHeaders)))

	t.Execute(&body, struct {
		Name    string
		Message string
	}{
		Name:    "Charlie martinez",
		Message: "This is a test message in a HTML template",
	})

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println("Email Sent!")

	return true
}
