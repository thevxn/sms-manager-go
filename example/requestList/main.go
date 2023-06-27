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
		Message:     "hello world!",
		GatewayType: "high",
	}

	resp := &client.Response{}

	if err := client.DoRequest(req, "RequestList", resp); err != nil {
		log.Fatal(err)
	}

	// print GetUserInfo response fields
	if resp.Message != "ERROR" {
		count := len(resp.RequestIDArray)

		for i := 0; i < count; i++ {
			fmt.Printf("request ID                 : %d\n\r", resp.RequestIDArray[i])
			fmt.Println("gateway type               : " + resp.GatewayTypeArray[i])
			fmt.Println("SMS sending time           : " + resp.TimeSentArray[i])
			fmt.Println("SMS send exp. time         : " + resp.TimeExpiryArray[i])
			fmt.Println("sender name set            : " + resp.SenderArray[i])
			fmt.Printf("remaining receivers count  : %d\n\r", resp.RemainCountArray[i])
			fmt.Printf("request process code/status: %d (%s)\n\r\n\r", resp.ProcessCodeArray[i], client.ProcessCodeText[resp.ProcessCodeArray[i]])
		}
	} else {
		fmt.Printf("response error code: %d\n\r", resp.ErrorCode)
		fmt.Println("response error msq : " + resp.ErrorMessage)
	}
}
