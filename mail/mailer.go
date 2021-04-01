package mail

import (
	"bytes"
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/tonsV2/event-rooster-api/configurations"
	"github.com/tonsV2/event-rooster-api/models"
	"gopkg.in/mail.v2"
	"html/template"
)

var testEmails = []string{"test@mail.com", "test1@mail.com", "test2@mail.com"}

func ProvideMailer(mailerConfiguration configurations.MailerConfiguration) Mailer {
	templatePathPrefix := "./"
	if flag.Lookup("test.v") != nil {
		templatePathPrefix = "../"
	}

	return Mailer{
		configuration: mailerConfiguration,

		from: "sebastianthegreatful@something.com",

		createEventSubject:  "New event created",
		createEventTemplate: templatePathPrefix + "mail/templates/createEvent.html",

		welcomeParticipantSubject:  "Welcome to the event",
		welcomeParticipantTemplate: templatePathPrefix + "mail/templates/welcomeEvent.html",
	}
}

type Mailer struct {
	configuration configurations.MailerConfiguration
	from          string

	createEventSubject  string
	createEventTemplate string

	welcomeParticipantSubject  string
	welcomeParticipantTemplate string
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

func (m *Mailer) SendWelcomeParticipantMail(event models.Event, participant models.Participant) error {
	t, _ := template.ParseFiles(m.welcomeParticipantTemplate)
	var body bytes.Buffer
	_ = t.Execute(&body, struct {
		EventId uint
		Title   string
		Message string
		Token   string
	}{
		EventId: event.ID,
		Title:   event.Title,
		Message: "This is a test message in a HTML template",
		Token:   participant.Token,
	})

	return m.sendMail(m.from, participant.Email, m.welcomeParticipantSubject, body)
}

func (m *Mailer) sendMail(from string, to string, subject string, body bytes.Buffer) error {
	message := mail.NewMessage()
	message.SetHeader("From", from)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", body.String())
	d := mail.NewDialer(m.configuration.Host, m.configuration.Port, m.configuration.Username, m.configuration.Password)
	if contains(testEmails, to) {
		return nil
	} else {
		d.TLSConfig = &tls.Config{InsecureSkipVerify: false, ServerName: m.configuration.Host}
		return d.DialAndSend(message)
	}
}

func (m *Mailer) SendWelcomeToGroupMail(event models.Event, group models.Group, participant models.Participant) error {
	fmt.Println("TODO: Send welcome to group mail")
	return nil
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
