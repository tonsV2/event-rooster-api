package mail

import (
	"bytes"
	"crypto/tls"
	"flag"
	"github.com/dustin/go-humanize"
	"github.com/tonsV2/event-rooster-api/configurations"
	"github.com/tonsV2/event-rooster-api/models"
	"gopkg.in/mail.v2"
	"html/template"
	"log"
	"strconv"
)

var testEmails = []string{"test@mail.com", "test1@mail.com", "test2@mail.com", "test3@mail.com", "test4@mail.com"}

func ProvideMailer(mailerConfiguration configurations.MailerConfiguration) Mailer {
	templatePathPrefix := "./"
	if flag.Lookup("test.v") != nil {
		templatePathPrefix = "../"
	}

	return Mailer{
		configuration: mailerConfiguration,

		createEventSubject:  "New event created",
		createEventTemplate: templatePathPrefix + "mail/templates/createEvent.html",

		welcomeParticipantSubject:  "Welcome to the event",
		welcomeParticipantTemplate: templatePathPrefix + "mail/templates/welcomeEvent.html",

		joinGroupSubject:  "Welcome to",
		joinGroupTemplate: templatePathPrefix + "mail/templates/joinGroup.html",
	}
}

type Mailer struct {
	configuration configurations.MailerConfiguration

	createEventSubject  string
	createEventTemplate string

	welcomeParticipantSubject  string
	welcomeParticipantTemplate string

	joinGroupSubject  string
	joinGroupTemplate string
}

func (m *Mailer) SendCreateEventMail(event models.Event) error {
	to := event.Email

	t, _ := template.ParseFiles(m.createEventTemplate)
	var body bytes.Buffer
	_ = t.Execute(&body, struct {
		DomainName string
		Title      string
		Token      string
	}{
		DomainName: m.configuration.DomainName,
		Title:      event.Title,
		Token:      event.Token,
	})

	subject := m.createEventSubject + ": " + event.Title
	return m.sendMail(m.configuration.Username, to, subject, body)
}

func (m *Mailer) SendWelcomeParticipantMail(event models.Event, participant models.Participant) error {
	return m.SendWelcomeParticipantMails(event, []models.Participant{participant})
}

func (m *Mailer) SendWelcomeParticipantMails(event models.Event, participants []models.Participant) error {
	var messages []*mail.Message
	for _, participant := range participants {
		body := m.generateBody(event, participant)
		message := m.generateMessage(m.configuration.Username, participant.Email, m.createEventSubject, body)
		messages = append(messages, message)
	}
	return m.sendMails(messages)
}

func (m *Mailer) generateBody(event models.Event, participant models.Participant) bytes.Buffer {
	t, _ := template.ParseFiles(m.welcomeParticipantTemplate)
	var body bytes.Buffer
	data := struct {
		DomainName string
		EventId    uint
		Title      string
		Token      string
	}{
		DomainName: m.configuration.DomainName,
		EventId:    event.ID,
		Title:      event.Title,
		Token:      participant.Token,
	}
	_ = t.Execute(&body, data)

	return body
}

func (m *Mailer) generateMessage(from string, to string, subject string, body bytes.Buffer) *mail.Message {
	message := mail.NewMessage()
	message.SetHeader("From", from)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", body.String())
	return message
}

// Inspiration: https://github.com/go-mail/mail/blob/v2.3.1/example_test.go#L77
func (m *Mailer) sendMails(messages []*mail.Message) error {
	d := mail.NewDialer(m.configuration.Host, m.configuration.Port, m.configuration.Username, m.configuration.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: false, ServerName: m.configuration.Host}

	s, err := d.Dial()
	if err != nil {
		return err
	}

	for _, message := range messages {
		log.Printf("Mail sent: %+v", message)
		to := message.GetHeader("To")[0]
		if contains(testEmails, to) {
			continue
		} else {
			err := mail.Send(s, message)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *Mailer) sendMail(from string, to string, subject string, body bytes.Buffer) error {
	message := mail.NewMessage()
	message.SetHeader("From", from)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", body.String())
	d := mail.NewDialer(m.configuration.Host, m.configuration.Port, m.configuration.Username, m.configuration.Password)
	log.Printf("Mail sent: %+v", message)
	if contains(testEmails, to) {
		return nil
	} else {
		d.TLSConfig = &tls.Config{InsecureSkipVerify: false, ServerName: m.configuration.Host}
		return d.DialAndSend(message)
	}
}

func (m *Mailer) SendWelcomeToGroupMail(event models.Event, group models.Group, participant models.Participant) error {
	to := event.Email

	gid, err := strconv.Atoi(group.GID)
	if err != nil {
		log.Fatal(err)
	}
	gidOrdinalized := humanize.Ordinal(gid)

	t, _ := template.ParseFiles(m.joinGroupTemplate)
	var body bytes.Buffer
	_ = t.Execute(&body, struct {
		DomainName string
		GID        string
		Token      string
		EventId    uint
	}{
		DomainName: m.configuration.DomainName,
		GID:        gidOrdinalized,
		Token:      participant.Token,
		EventId:    event.ID,
	})

	subject := m.joinGroupSubject + " " + gidOrdinalized + " group"
	return m.sendMail(m.configuration.Username, to, subject, body)
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
