package email

import (
	"github.com/kunaldawn/mail.af/pkg/db/models"
	"gopkg.in/gomail.v2"
)

type Sender struct {
	sender *models.Sender
	dialer *gomail.Dialer
	closer gomail.SendCloser
}

func NewSender(sender *models.Sender) (*Sender, error) {
	context := &Sender{sender: sender, dialer: gomail.NewPlainDialer("smtp.gmail.com", 587, sender.Email, sender.Password)}
	closer, err := context.dialer.Dial()
	if err != nil {
		return &Sender{}, nil
	}
	context.closer = closer

	return context, nil
}

func (sender *Sender) Send(from string, to string, message *gomail.Message) error {
	return sender.closer.Send(from, []string{to}, message)
}
