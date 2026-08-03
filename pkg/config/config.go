package config

import (
	"flag"
)

var (
	APIKey string
	Method string

	PhoneNumber string
	Body        string
	Sender      string
	Gateway     string
	Tag         string
	DateTime    string

	MessageID string
	RequestID string
	Date      string

	PushGatewayURL string
)

func init() {
	apiKey := flag.String("apiKey", "", "a string, API key issued by the SMS Manager portal")
	method := flag.String("method", "", "a string, the client method name")

	phoneNumber := flag.String("phoneNumber", "", "a string, a phone number to send the message to")
	body := flag.String("body", "", "a string, the message text")
	sender := flag.String("sender", "", "a string, the sender name or virtual number")
	gateway := flag.String("gateway", "high", "a string, the SMS gateway used --- high, lowcost, direct, custom, simhost or gsm")
	tag := flag.String("tag", "", "a string, the message tag --- promotional, transactional, priority, or an own grouping value")
	dateTime := flag.String("dateTime", "", "a string, an optional scheduled sending time in UTC; format: 2023-12-31T23:59:59Z")

	messageID := flag.String("messageID", "", "a string, an identificator of a single sent message")
	requestID := flag.String("requestID", "", "a string, an identificator of a sending request")
	date := flag.String("date", "", "a string, a date to list the messages/statistics for; format: 2023-12-31")

	pushGatewayURL := flag.String("pushGatewayURL", "", "a string, an optional URL target for fetched values reporting as metrics")

	flag.Parse()

	// compulsory vars
	APIKey = *apiKey
	Method = *method

	// sending vars
	PhoneNumber = *phoneNumber
	Body = *body
	Sender = *sender
	Gateway = *gateway
	Tag = *tag
	DateTime = *dateTime

	// status lookup vars
	MessageID = *messageID
	RequestID = *requestID
	Date = *date

	// very optional vars
	PushGatewayURL = *pushGatewayURL
}
