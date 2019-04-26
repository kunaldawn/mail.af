package models

import "time"

type Log struct {
	JobID    string    `json:"job_id" bson:"job_id"`
	Receiver *Receiver `json:"receiver" bson:"receiver"`
	Time     int64     `json:"time" bson:"time"`
	Success  bool      `json:"success" bson:"success"`
}

func NewLog(id string, receiver *Receiver, success bool) *Log {
	return &Log{JobID: id, Receiver: receiver, Time: time.Now().UnixNano(), Success: success}
}
