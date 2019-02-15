package notification

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"
	"time"
)

type mail struct {
	sender  string
	to      []string
	cc      []string
	subject string
	body    string
}

type SmtpConfig struct {
	Host     string
	Port     string
	User     string
	Password string
}

func (mail *mail) buildMessage() string {
	msg := ""
	msg += fmt.Sprintf("From: %s\r\n", mail.sender)
	if len(mail.to) > 0 {
		msg += fmt.Sprintf("To: %s\r\n", strings.Join(mail.to, ";"))
	}
	if len(mail.cc) > 0 {
		msg += fmt.Sprintf("Cc: %s\r\n", strings.Join(mail.cc, ";"))
	}

	msg += fmt.Sprintf("Subject: %s\r\n", mail.subject)
	msg += fmt.Sprintf("Date: %v\r\n", time.Now())
	msg += "\r\n" + mail.body

	return msg
}

func (config *SmtpConfig) buildServerName() string {
	return fmt.Sprintf("%s:%s", config.Host, config.Port)
}

func (config *SmtpConfig) sendMail(mail mail) error {
	auth := smtp.PlainAuth("", config.User, config.Password, config.Host)

	msg := mail.buildMessage()
	err := smtp.SendMail(config.buildServerName(), auth, mail.sender, mail.to, []byte(msg))
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
