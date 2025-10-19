package utils

import (
	"fmt"

	"github.com/AbdulRahman-04/07EvenetManagement-API/server/config"
	"gopkg.in/gomail.v2"
)

type EmailData struct{
	From string
	To string
	Subject string
	Text string
	HTML string
}

func SendEmail(data EmailData) {

	// get user nad pass
	User := config.AppConfig.Email.USER
	Pass := config.AppConfig.Email.PASS

	// create sender 
	s := gomail.NewMessage()

	s.SetAddressHeader("From", User, "BookMyEvents.com")
	s.SetHeader("To", data.To)
	s.SetHeader("Subject", data.Subject)
	s.SetBody("text/plain", data.Text)
	s.AddAlternative("text/html", data.HTML)

	// create transporter 
	t := gomail.NewDialer("smtp.gmail.com", 465, User, Pass)

	// try sending mail
	err := t.DialAndSend(s)
	if err != nil {
		fmt.Println("error sending email")
		return
	}
	fmt.Println("Email sent✅")
}