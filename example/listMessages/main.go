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

	params := client.ListMessagesParams{
		Date: time.Now().Format("2006-01-02"),

		// optional filters
		//Tag:         "newsletter",
		//PhoneNumber: "420123456789",
	}

	// page through the whole day
	for {
		list, err := c.ListMessages(ctx, params)
		if err != nil {
			log.Fatal(err)
		}

		for _, msg := range list.Messages {
			fmt.Println("message ID     : " + msg.MessageID)
			fmt.Println("phone number   : " + msg.PhoneNumber)
			fmt.Println("channel used   : " + msg.Flow)
			fmt.Println("sender name set: " + msg.Sender)
			fmt.Printf("message parts  : %d\n\r", msg.MessageParts)
			fmt.Printf("sent/delivered : %d/%d\n\r", msg.Sent, msg.Delivered)
			fmt.Printf("price          : %s CZK\n\r\n\r", msg.PriceCZK)
		}

		if list.NextCursor == nil {
			break
		}

		params.Cursor = *list.NextCursor
	}
}
