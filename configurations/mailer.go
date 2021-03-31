package configurations

import (
	"strconv"
)

func ProvideMailerConfiguration() MailerConfiguration {
	host := requireEnv("SMTP_HOST")
	portStr := requireEnv("SMTP_PORT")
	port, _ := strconv.Atoi(portStr)
	username := requireEnv("SMTP_USERNAME")
	password := requireEnv("SMTP_PASSWORD")

	return MailerConfiguration{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}
}

type MailerConfiguration struct {
	Host     string
	Port     int
	Username string
	Password string
}
