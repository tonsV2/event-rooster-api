package mail

import (
	"bytes"
	"crypto/tls"
	"github.com/tonsV2/event-rooster-api/configurations"
	"github.com/tonsV2/event-rooster-api/models"
	"gopkg.in/mail.v2"
	"html/template"
)

func ProvideMailer(mailerConfiguration configurations.MailerConfiguration) Mailer {
	return Mailer{
		configuration: mailerConfiguration,

		from: "sebastianthegreatful@something.com",

		createEventSubject:  "New event created",
		createEventTemplate: "./mail/templates/createEvent.html",

		welcomeRunnerSubject:  "Welcome to the event",
		welcomeRunnerTemplate: "./mail/templates/welcomeEvent.html",
	}
}

type Mailer struct {
	configuration configurations.MailerConfiguration
	from          string

	createEventSubject  string
	createEventTemplate string

	welcomeRunnerSubject  string
	welcomeRunnerTemplate string
}

func (m *Mailer) SendCreateEventMail(event models.Event) error {
	to := event.Email

	t, _ := template.ParseFiles(m.createEventTemplate)
	var body bytes.Buffer
	_ = t.Execute(&body, struct {
		Title   string
		Message string
		Token   string
	}{
		Title:   event.Title,
		Message: "This is a test message in a HTML template",
		Token:   event.Token,
	})

	return m.sendMail(m.from, to, m.createEventSubject, body)
}

func (m *Mailer) SendWelcomeRunnerMail(event models.Event, runner models.Runner) error {
	to := event.Email

	t, _ := template.ParseFiles(m.welcomeRunnerTemplate)
	var body bytes.Buffer
	_ = t.Execute(&body, struct {
		Title   string
		Message string
		Token   string
	}{
		Title:   event.Title,
		Message: "This is a test message in a HTML template",
		Token:   event.Token,
	})

	return m.sendMail(m.from, to, m.welcomeRunnerSubject, body)
}

func (m *Mailer) sendMail(from string, to string, subject string, body bytes.Buffer) error {
	message := mail.NewMessage()
	message.SetHeader("From", from)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", body.String())
	d := mail.NewDialer(m.configuration.Host, m.configuration.Port, m.configuration.Username, m.configuration.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: false, ServerName: m.configuration.Host}
	return d.DialAndSend(message)
}
