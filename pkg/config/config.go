package config

import (
	"flag"
)

var (
	MethodName string
	APIKey string
	PhoneNumber string
	GatewayType string

	RequestID string
	Sender string
	CustomID string
	SendTime string
	ExpirationTime string
)

func init() {
	methodName := flag.String("method", "", "a string, the client method name")
	apiKey := flag.String("apiKey", "", "a string, API key issued by SMS Manager portal")
	phoneNumber := flag.String("phoneNumber", "", "a string, a phone number to send the SMS to")
	gatewayType := flag.String("high", "a string, the specification of gateway service/tariff used")
	
	requestID := flag.String("requestID", "", "a string, an identificator returned by the service on sent SMS")
	sender := flag.String("sender", "info-sms", "a string, the SMS sender identificator")
	customID := flag.String("customID", "", "a string, a custom identificator for such SMS being sent")
	sendTime := flag.String("sendTime", "", "a string, an optional specificaton of the SMS sending time; format: 2023-12-32T23:59:59")
	expirationTime := flag.String("expirationTime", "", "a string, an optional specificaton of the SMS expiration time; format, the SMS won't be sent after that datetime specified: 2023-12-32T23:59:59")

	flag.Parse()

	// compulsory vars
	MethodName = *methodName
	APIKey = *apiKey
	PhoneNumber = *phoneNumber
	GatewayType = *gatewayType

	// optional vars
	RequestID = *requestID
	Sender = *sender
	CustomID = *customID
	SendTime = *sendTime
	ExpirationTime = *expirationTime
}


