package configurations

import (
	"strconv"
)

func ProvideMailerConfiguration() MailerConfiguration {
	domainName := requireEnv("DOMAIN_NAME")

	host := requireEnv("SMTP_HOST")
	portStr := requireEnv("SMTP_PORT")
	port, _ := strconv.Atoi(portStr)
	username := requireEnv("SMTP_USERNAME")
	password := requireEnv("SMTP_PASSWORD")

	return MailerConfiguration{
		DomainName: domainName,
		Host:       host,
		Port:       port,
		Username:   username,
		Password:   password,
	}
}

type MailerConfiguration struct {
	DomainName string
	Host       string
	Port       int
	Username   string
	Password   string
}
