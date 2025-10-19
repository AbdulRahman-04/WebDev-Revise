package utils

import (
	"fmt"

	"github.com/AbdulRahman-04/07EvenetManagement-API/server/config"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type SMSdata struct {
	From string
	To string
	Body string
}

func SendSMS(data SMSdata){
	// create client 
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username : config.AppConfig.Phone.SID,
		Password: config.AppConfig.Phone.TOKEN,
	})

	// get body ready 
	_, err := client.Api.CreateMessage(&openapi.CreateMessageParams{
		From: &config.AppConfig.Phone.PHONE,
		To: &data.To,
		Body: &data.Body,
	})

	if err != nil {
		fmt.Println("err sending mail❌")
		return
	}

	fmt.Println("SMS sent✅")
}