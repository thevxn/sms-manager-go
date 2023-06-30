package main

import (
	"fmt"
	"log"

	client "go.savla.dev/sms-manager/pkg/client"
)

func main() {
	req := client.Request{
		// compulsory params
		APIKey:      "xyz",
		PhoneNumber: "420123456789",
		Message:     "hello world!",
		GatewayType: "high",
	}

	resp := &client.Response{}

	if err := client.DoRequest(req, "GetUserInfo", resp); err != nil {
		log.Fatal(err)
	}

	// print GetUserInfo response fields
	if resp.Message != "ERROR" {
		fmt.Printf("credit left : %.2f CZK (exc. VAT)\n\r", resp.CreditBalance)
		fmt.Println("sender name : " + resp.Sender)
		fmt.Println("gateway type: " + resp.GatewayType)
	} else {
		fmt.Printf("response error code: %d\n\r", resp.ErrorCode)
		fmt.Println("response error msq : " + resp.ErrorMessage)
	}
}
