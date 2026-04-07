package mailer

import (
    "github.com/go-mail/mail/v2"
)

type Mailer struct {
    dialer *mail.Dialer
    sender string
}

func New(host string, port int, username, password, sender string) Mailer {
    dialer := mail.NewDialer(host, port, username, password)

    return Mailer{
        dialer: dialer,
        sender: sender,
    }
}

func (m Mailer) Send(to string, subject string, body string) error {
    msg := mail.NewMessage()

    msg.SetHeader("To", to)
    msg.SetHeader("From", m.sender)
    msg.SetHeader("Subject", subject)

    msg.SetBody("text/plain", body)

    return m.dialer.DialAndSend(msg)
}