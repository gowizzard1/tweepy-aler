package emailNotification

import (
	"net/smtp"
	"strings"
)

// Config contains the email and SMTP server configuration
type Config struct {
	ConsumerKey    string
	ConsumerSecret string
	AccessToken    string
	AccessSecret   string
	SMTPServer     string
	SMTPPort       string
	SMTPEmail      string
	SMTPPassword   string
}

// Email contains the information for an email notification
type Email struct {
	SenderID string
	ToIDs    []string
	Subject  string
	Message  string
}

// SendEmail sends an email notification
func SendEmail(config *Config, email *Email) error {
	auth := smtp.PlainAuth("", config.SMTPEmail, config.SMTPPassword, config.SMTPServer)
	to := strings.Join(email.ToIDs, ",")
	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + email.Subject + "\r\n" +
		"\r\n" +
		email.Message + "\r\n")
	return smtp.SendMail(config.SMTPServer+":"+config.SMTPPort, auth, config.SMTPEmail, email.ToIDs, msg)
}
