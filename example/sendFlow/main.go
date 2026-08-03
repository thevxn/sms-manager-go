package main

import (
	"context"
	"fmt"
	"log"
	"time"

	client "go.vxn.dev/sms-manager/pkg/client"
)

func main() {
	c, err := client.New("xyz")
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// the message is sent over Viber first, if it is not delivered within 10
	// minutes, it is sent again as an SMS
	resp, err := c.Send(ctx, client.Message{
		To: client.Recipients("420123456789"),
		Flow: []client.FlowStep{
			{
				Viber: &client.ViberStep{
					Body:         "your order is ready for pickup",
					Sender:       "MyCompany",
					TTL:          10,
					TTLCondition: client.TTLConditionDelivered,
				},
			},
			{
				SMS: &client.SMSStep{
					Body:    "your order is ready for pickup",
					Sender:  "MyCompany",
					Gateway: client.GatewayHigh,
					Type:    client.SMSTypeUTF,
				},
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("request ID: " + resp.RequestID)

	for _, accepted := range resp.Accepted {
		fmt.Printf("accepted recipient %s: message ID %s\n\r", accepted.Key, accepted.MessageID)
	}
}
