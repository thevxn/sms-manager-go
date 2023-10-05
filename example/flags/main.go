package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

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
		os.Exit(1)
	}

	// print method's response fields --- only raw credit balance there for example
	if resp.Message != "ERROR" {
		//fmt.Printf("credit left : %.2f CZK (exc. VAT)\n\r", resp.CreditBalance)
		//fmt.Printf("%.2f", resp.CreditBalance)
	} else {
		fmt.Printf("response error code: %d\n\r", resp.ErrorCode)
		fmt.Println("response error msq : " + resp.ErrorMessage)
		os.Exit(2)
	}

	// compose a request body
	body := fmt.Sprintln("#HELP SMS Manager available credit balance")
	body += fmt.Sprintln("#TYPE credit_balance_czk_total gauge")
	body += fmt.Sprintln("credit_balance_czk_total", strconv.FormatFloat(resp.CreditBalance, 'f', 2, 64))

	bodyReader := bytes.NewReader([]byte(body))

	// report the balance to PushGateway
	request, err := http.NewRequest("PUT", config.PushGatewayURL, bodyReader)
	if err != nil {
		log.Fatal(err)
		os.Exit(3)
	}
	request.Header.Set("Content-Type", "application/json")

	client := http.Client{}

	res, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
		os.Exit(4)
	}

	log.Println("Results pushed to pushgateway")
	defer res.Body.Close()
}

