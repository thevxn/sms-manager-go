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

	resp, err := c.Send(ctx, client.Message{
		Body:   "hello world!",
		To:     client.Recipients("420123456789"),
		Sender: "SmsManager",
		Tag:    client.TagTransactional,

		// optional metadata returned back in delivery notifications
		Payload: map[string]interface{}{"customID": "123456"},

		// optional scheduling
		//DateTime: client.NewTime(time.Now().Add(time.Hour)),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("request ID: " + resp.RequestID)

	for _, accepted := range resp.Accepted {
		fmt.Printf("accepted recipient %s: message ID %s\n\r", accepted.Key, accepted.MessageID)
	}

	for _, rejected := range resp.Rejected {
		fmt.Printf("rejected recipient %s\n\r", rejected.Key)
	}
}
