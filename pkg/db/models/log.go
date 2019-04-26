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
