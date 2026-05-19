package provider

import (
	"fmt"
	"net/smtp"
)

type SMTPEmailSender struct {
	host     string
	port     string
	username string
	password string
	from     string
}

func NewSMTPEmailSender(host, port, username, password, from string) *SMTPEmailSender {
	return &SMTPEmailSender{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

func (s *SMTPEmailSender) Send(to, subject, body string) error {
	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	msg := []byte(fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		s.from, to, subject, body,
	))

	return smtp.SendMail(addr, auth, s.from, []string{to}, msg)
}
