// Filename: test2/internal/mailer/mailer.go
package mailer

import (
	"bytes"
	"embed"
	"text/template"
	"time"

	"gopkg.in/mail.v2"
)

// go:embed "templates"
var tempalteFS embed.FS

// creating the mailer type
type Mailer struct {
	dailer *mail.Dialer
	sender string
}

// New() creates a new mailer object
func New(host string, port int, username, password, sender string) Mailer {
	dailer := mail.NewDialer(host, port, username, password)
	dailer.Timeout = 5 * time.Second

	//returning the instance of the mailer
	return Mailer{
		dailer: dailer,
		sender: sender,
	}
}

// Sending the mail
func (m Mailer) Send(recipient, templateFile string, data interface{}) error {
	tmpl, err := template.New("email").ParseFS(tempalteFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	//Executing the tempate
	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	//Executing the template again
	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}

	//Executing the template again again
	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	//Creating the new mail message
	msg := mail.NewMessage()
	msg.SetHeader("To", recipient)
	msg.SetHeader("From", m.sender)
	msg.SetHeader("Subject", subject.String())
	msg.SetBody("text/plain", plainBody.String())
	msg.AddAlternative("text/html", htmlBody.String())

	//Calling the Dail and Send functions
	err = m.dailer.DialAndSend(msg)
	if err != nil {
		return err
	}
	return nil
}
