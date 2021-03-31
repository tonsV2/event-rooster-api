package mail

import (
	"bytes"
	"crypto/tls"
	"github.com/tonsV2/race-rooster-api/configurations"
	"github.com/tonsV2/race-rooster-api/models"
	"gopkg.in/mail.v2"
	"html/template"
)

func ProvideMailer(mailerConfiguration configurations.MailerConfiguration) Mailer {
	return Mailer{
		configuration: mailerConfiguration,

		from: "sebastianthegreatful@something.com",

		createRaceSubject:  "New race created",
		createRaceTemplate: "./mail/templates/createRace.html",

		welcomeRunnerSubject:  "Welcome to the race",
		welcomeRunnerTemplate: "./mail/templates/welcomeRunner.html",
	}
}

type Mailer struct {
	configuration configurations.MailerConfiguration
	from          string

	createRaceSubject  string
	createRaceTemplate string

	welcomeRunnerSubject  string
	welcomeRunnerTemplate string
}

func (m *Mailer) SendCreateRaceMail(race models.Race) error {
	to := race.Email

	t, _ := template.ParseFiles(m.createRaceTemplate)
	var body bytes.Buffer
	_ = t.Execute(&body, struct {
		Title   string
		Message string
		Token   string
	}{
		Title:   race.Title,
		Message: "This is a test message in a HTML template",
		Token:   race.Token,
	})

	return m.sendMail(m.from, to, m.createRaceSubject, body)
}

func (m *Mailer) SendWelcomeRunnerMail(race models.Race, runner models.Runner) error {
	to := race.Email

	t, _ := template.ParseFiles(m.welcomeRunnerTemplate)
	var body bytes.Buffer
	_ = t.Execute(&body, struct {
		Title   string
		Message string
		Token   string
	}{
		Title:   race.Title,
		Message: "This is a test message in a HTML template",
		Token:   race.Token,
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
