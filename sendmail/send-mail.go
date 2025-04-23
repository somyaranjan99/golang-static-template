package sendmail

import (
	"fmt"
	"github/somyaranjan99/basic-go-project/pkg/config"
	"github/somyaranjan99/basic-go-project/pkg/models"
	"log"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
)

type NewMail struct {
	MailDetails chan models.MailData
}

func NewMAilRepo(a *config.AppConfig) *NewMail {
	return &NewMail{MailDetails: a.MailChan}
}
func (n *NewMail) ListenForMail() {
	go func() {
		msg := <-n.MailDetails
		sendMsg(msg)
	}()
}

func sendMsg(m models.MailData) {
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second
	client, err := server.Connect()
	if err != nil {
		fmt.Println(err)
	}
	email := mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)
	email.SetBody(mail.TextHTML, "Hello,<strong>World</strong>")
	err = email.Send(client)
	if err != nil {
		log.Println(err)
	} else {
		println("Email Sent")
	}
}
