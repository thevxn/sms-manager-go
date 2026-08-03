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

	prices, err := c.Prices(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// price list is keyed by channel and by the country calling code
	price, ok := prices["sms"]["420"]
	if !ok {
		log.Fatal("no SMS price for the 420 calling code")
	}

	fmt.Printf("destination MCC: %d\n\r", price.MCC)
	fmt.Printf("price per SMS  : %.2f CZK / %.3f EUR\n\r", price.PriceCZK, price.PriceEUR)

	for operator, operatorPrice := range price.Operators {
		fmt.Printf("operator %s     : %.2f CZK / %.3f EUR\n\r", operator, operatorPrice.PriceCZK, operatorPrice.PriceEUR)
	}
}
