/*
 __  __       _ _      _    _____
|  \/  | __ _(_) |    / \  |  ___|
| |\/| |/ _` | | |   / _ \ | |_
| |  | | (_| | | |_ / ___ \|  _|
|_|  |_|\__,_|_|_(_)_/   \_\_|

Send mails as fuck!
Author : Kunal Dawn (kunal.dawn@gmail.com)

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>
*/
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
