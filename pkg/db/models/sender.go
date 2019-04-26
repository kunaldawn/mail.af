package models

import "github.com/kunaldawn/mail.af/pkg/utils"

type Sender struct {
	ID       string `json:"id" bson:"id"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

func NewSender(email string, password string) *Sender {
	return &Sender{ID: utils.UUID(), Email: email, Password: password}
}
