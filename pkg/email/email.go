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
package email

import (
	"encoding/base64"
	"fmt"
	"github.com/kunaldawn/mail.af/pkg/db/models"
	"gopkg.in/gomail.v2"
	"io"
	"strings"
)

func Send(sender *Sender, receiver *models.Receiver, subject, image string) error {
	message := gomail.NewMessage()
	message.SetHeader("From", sender.sender.Email)
	message.SetHeader("To", receiver.Email)
	message.SetHeader("Subject", subject)

	tagData := image[:strings.IndexByte(image, ',')]
	b64Data := image[strings.IndexByte(image, ',')+1:]
	imageName := "image.png"
	if strings.Contains(tagData, "jpeg") || strings.Contains(tagData, "jpg") {
		imageName = "image.jpeg"
	}

	message.SetBody("text/html", fmt.Sprintf("<img src=\"cid:%s\"/>", imageName))
	message.Embed(imageName, gomail.SetCopyFunc(func(writer io.Writer) error {
		if data, err := base64.StdEncoding.DecodeString(b64Data); err == nil {
			writer.Write(data)
		}
		return nil
	}))

	return sender.Send(sender.sender.Email, receiver.Email, message)
}
