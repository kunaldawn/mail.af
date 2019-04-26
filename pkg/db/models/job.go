package models

import "github.com/kunaldawn/mail.af/pkg/utils"

type Job struct {
	ID          string      `json:"id" bson:"id"`
	Name        string      `json:"name" bson:"name"`
	StartTime   int64       `json:"start_time" bson:"start_time"`
	EndTime     int64       `json:"end_time" bson:"end_time"`
	Done        bool        `json:"done" bson:"done"`
	Running     bool        `json:"running" bson:"running"`
	Subject     string      `json:"subject" bson:"subject"`
	Image       string      `json:"image" bson:"image"`
	Receivers   []*Receiver `json:"receivers" bson:"receivers"`
	SendSuccess int         `json:"send_success" bson:"send_success"`
	SendFailed  int         `json:"send_failed" bson:"send_failed"`
}

func NewJob(name string, start int64, subject string, image string, receivers []*Receiver) *Job {
	return &Job{ID: utils.UUID(), Name: name, StartTime: start, Subject: subject, Image: image, Receivers: receivers}
}
