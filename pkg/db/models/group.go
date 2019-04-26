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
