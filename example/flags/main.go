package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	client "go.vxn.dev/sms-manager/pkg/client"
	config "go.vxn.dev/sms-manager/pkg/config"
)

func main() {
	c, err := client.New(config.APIKey)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	switch config.Method {
	case "send":
		msg := client.Message{
			Body:   config.Body,
			To:     client.Recipients(config.PhoneNumber),
			Sender: config.Sender,
			Tag:    config.Tag,
			Flow: []client.FlowStep{
				{SMS: &client.SMSStep{Gateway: client.Gateway(config.Gateway)}},
			},
		}

		if config.DateTime != "" {
			sendTime, err := time.Parse(time.RFC3339, config.DateTime)
			if err != nil {
				log.Fatal(err)
			}
			msg.DateTime = client.NewTime(sendTime)
		}

		resp, err := c.Send(ctx, msg)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("request ID: " + resp.RequestID)

		for _, accepted := range resp.Accepted {
			fmt.Println("message ID: " + accepted.MessageID)
		}

		for _, rejected := range resp.Rejected {
			fmt.Printf("rejected recipient: %s\n\r", rejected.Key)
			os.Exit(2)
		}

	case "credit":
		credit, err := c.Credit(ctx)
		if err != nil {
			log.Fatal(err)
		}

		balance, err := credit.Balance()
		if err != nil {
			log.Fatal(err)
		}

		if config.PushGatewayURL == "" {
			fmt.Printf("%.2f\n\r", balance)
			return
		}

		if err := pushCreditBalance(config.PushGatewayURL, balance); err != nil {
			log.Fatal(err)
		}

		log.Println("Results pushed to pushgateway")

	default:
		log.Fatal("unknown method, use -method send or -method credit")
	}
}

// pushCreditBalance reports the credit balance to a Prometheus PushGateway.
func pushCreditBalance(pushGatewayURL string, balance float64) error {
	body := fmt.Sprintln("#HELP SMS Manager available credit balance")
	body += fmt.Sprintln("#TYPE credit_balance_czk_total gauge")
	body += fmt.Sprintln("credit_balance_czk_total", strconv.FormatFloat(balance, 'f', 2, 64))

	request, err := http.NewRequest(http.MethodPut, pushGatewayURL, bytes.NewReader([]byte(body)))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "text/plain")

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
