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
