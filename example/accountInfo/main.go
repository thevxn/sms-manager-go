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

	info, err := c.APIKeyInfo(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("account type : " + info.Account)
	fmt.Println("key scope    : " + info.Scope)
	fmt.Println("currency     : " + info.Currency)

	for _, sender := range info.Senders {
		fmt.Printf("sender       : %s (%s)\n\r", sender.Value, sender.Flow)
	}

	credit, err := c.Credit(ctx)
	if err != nil {
		log.Fatal(err)
	}

	balance, err := credit.Balance()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("credit left  : %.2f %s (exc. VAT)\n\r", balance, credit.Currency)
}
