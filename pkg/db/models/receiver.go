package models

import "github.com/kunaldawn/mail.af/pkg/utils"

type Receiver struct {
	ID    string `json:"id" bson:"id"`
	Email string `json:"email" bson:"email"`
}

func NewReceiver(email string) *Receiver {
	return &Receiver{ID: utils.UUID(), Email: email}
}
