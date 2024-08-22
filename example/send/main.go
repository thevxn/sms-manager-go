package main

import (
	"fmt"
	"log"

	client "go.vxn.dev/sms-manager/pkg/client"
)

func main() {
	req := client.Request{
		// compulsory params
		APIKey:      "xyz",
		PhoneNumber: "420123456789",
		Message:     "hello world!",
		GatewayType: "high",
		// optional params
		CustomID: "123456",
		Sender:   "info-sms",
		//SendTime:       "2023-06-26T20:45:00",
		//ExpirationTime: "2023-06-26T20:55:00",
	}

	resp := &client.Response{}

	if err := client.DoRequest(req, "Send", resp); err != nil {
		log.Fatal(err)
	}

	if resp.Message != "ERROR" {
		fmt.Println("status         : " + resp.Message)
		fmt.Printf("request ID     : %d\n\r", resp.RequestID)
		fmt.Println("phone number(s): " + resp.PhoneNumber)
		fmt.Printf("custom ID      : %d\n\r", resp.CustomID)
	} else {
		fmt.Printf("response error code: %d\n\r", resp.ErrorCode)
		fmt.Println("response error msq : " + resp.ErrorMessage)
	}
}
