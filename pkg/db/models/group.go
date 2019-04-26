package models

import (
	"github.com/kunaldawn/mail.af/pkg/utils"
)

type Group struct {
	ID        string      `json:"id" bson:"id"`
	Name      string      `json:"name" bson:"name"`
	Receivers []*Receiver `json:"receivers" bson:"receivers"`
}

func NewGroup(name string, receivers []*Receiver) *Group {
	return &Group{ID: utils.UUID(), Name: name, Receivers: receivers}
}
