package main

import (
	"context"
	"errors"
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

	status, err := c.MessageStatus(ctx, "pppppppp-qqqq-rrrr-ssss-tttttttttttt")
	if errors.Is(err, client.ErrNotFound) {
		log.Fatal("no such message sent by this API key")
	}
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("phone number   : " + status.To.PhoneNumber)
	fmt.Println("send status    : " + status.StatusCode.Text())
	fmt.Println("delivery status: " + status.DeliveryCode.Text())

	if code, description := client.ParseStatusInfo(status.StatusInfo); code != "" {
		fmt.Printf("status info    : %s (%s)\n\r", description, client.ResultCodeText[code])
	}

	if status.TimeDelivered != nil {
		fmt.Println("delivered at   : " + status.TimeDelivered.Format(time.RFC3339))
	}

	if status.Price != nil {
		fmt.Printf("price          : %.2f\n\r", *status.Price)
	}
}
