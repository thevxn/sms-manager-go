package main

import (
	"fmt"
	"log"

	client "go.savla.dev/sms-manager/pkg/client"
	config "go.savla.dev/sms-manager/pkg/config"
)

func main() {
	// assign flags to all fields
	req := client.Request{
		// compulsory params
		APIKey:         config.APIKey,
		PhoneNumber:    config.PhoneNumber,
		Message:        "",
		GatewayType:    config.GatewayType,

		RequestID:      config.RequestID,
		Sender:         config.Sender,
		CustomID:       config.CustomID,
		SendTime:       config.SendTime,
		ExpirationTime: config.ExpirationTime,
	}

	resp := &client.Response{}

	// do the request
	if err := client.DoRequest(req, config.MethodName, resp); err != nil {
		log.Fatal(err)
	}

	// print method's response fields --- only raw credit balance there for example
	if resp.Message != "ERROR" {
		//fmt.Printf("credit left : %.2f CZK (exc. VAT)\n\r", resp.CreditBalance)
		fmt.Printf("%.2f", resp.CreditBalance)
	} else {
		fmt.Printf("response error code: %d\n\r", resp.ErrorCode)
		fmt.Println("response error msq : " + resp.ErrorMessage)
	}
}
