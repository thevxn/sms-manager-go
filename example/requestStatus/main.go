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
		RequestID: "58019737",
		CustomID:  "123456",
	}

	resp := &client.Response{}

	if err := client.DoRequest(req, "RequestStatus", resp); err != nil {
		log.Fatal(err)
	}

	// print GetUserInfo response fields
	if resp.Message != "ERROR" {
		fmt.Println("phone number      : " + resp.PhoneNumber)
		fmt.Printf("send status       : %d\n\r", resp.ProcessCode)
		fmt.Printf("delivery msg count: %d\n\r", resp.DeliveryCount)
		fmt.Printf("delivery status   : %d\n\r", resp.DeliveryCode)
	} else {
		fmt.Printf("response error code: %d\n\r", resp.ErrorCode)
		fmt.Println("response error msg: " + resp.ErrorMessage)
	}
}
