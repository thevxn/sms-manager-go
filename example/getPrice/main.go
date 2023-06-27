package main

import (
	"fmt"
	"log"

	client "go.savla.dev/sms-manager/pkg/client"
	config "go.savla.dev/sms-manager/pkg/config"
)

func main() {
	req := config.Request{
		// compulsory params
		APIKey:      "xyz",
		PhoneNumber: "420123456789",
		Message:     "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.",
		GatewayType: "high",
	}

	resp := &client.Response{}

	if err := client.DoRequest(req, "GetPrice", resp); err != nil {
		log.Fatal(err)
	}

	// print GetUserInfo response fields
	if resp.Message != "ERROR" {
		fmt.Printf("valid receivers count    : %d\n\r", resp.ValidReceiverCount)
		fmt.Printf("number of SMS msgs       : %d\n\r", resp.TextSplitSMSCount)
		fmt.Printf("SMS msg's character count: %d\n\r", resp.CharacterCount)
		fmt.Printf("price per one SMS msg    : %.2f CZK\n\r", resp.OneSMSPrice)
		fmt.Printf("price per all SMS msgs   : %.2f CZK\n\r", resp.SumSMSPrice)
	} else {
		fmt.Printf("response error code: %d\n\r", resp.ErrorCode)
		fmt.Println("response error msq : " + resp.ErrorMessage)
	}
}
