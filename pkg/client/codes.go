package client

import (
	"strconv"
	"strings"
)

// StatusCode is a message sending/delivery status code returned by the REST API.
type StatusCode string

const (
	// StatusSent is the sending status code of a sent message.
	StatusSent StatusCode = "200"

	// DeliveryDelivered is the delivery status code of a delivered message.
	DeliveryDelivered StatusCode = "201"

	// DeliverySeen is the delivery status code of a message seen by the recipient.
	DeliverySeen StatusCode = "202"
)

// StatusCodeText maps the known status codes to a human-readable description.
var StatusCodeText = map[StatusCode]string{
	StatusSent:        "sent",
	DeliveryDelivered: "delivered",
	DeliverySeen:      "seen",
}

// Text returns a human-readable description of the status code, or the raw code
// when it is not a known one.
func (c StatusCode) Text() string {
	if text, ok := StatusCodeText[c]; ok {
		return text
	}

	return string(c)
}

// Message results reported by the sentMessage webhook. A message is first
// reported as ResultSent, further states follow. ResultRejected means the
// message was not sent at all.
const (
	ResultSending     = "sending"
	ResultSent        = "sent"
	ResultDelivered   = "delivered"
	ResultUndelivered = "undelivered"
	ResultSeen        = "seen"
	ResultRejected    = "rejected"
	ResultFailed      = "failed"
)

// ResultCodeText maps the common rejection codes found in the status/result
// info strings to a human-readable description.
var ResultCodeText = map[string]string{
	"302": "invalid phone number format",
	"303": "phone number not in the allowed list",
	"304": "phone number blacklisted or opted out",
	"306": "too many identical messages to the same recipient",
	"307": "insufficient credit",
}

// ParseStatusInfo splits a status/result info string in the "[code] description"
// format into its parts. Both parts are optional --- an info string can carry a
// bracketed code only, or a description only.
func ParseStatusInfo(info string) (code string, description string) {
	info = strings.TrimSpace(info)

	if !strings.HasPrefix(info, "[") {
		return "", info
	}

	end := strings.Index(info, "]")
	if end < 0 {
		return "", info
	}

	return info[1:end], strings.TrimSpace(info[end+1:])
}

// Amount is a numeric value the REST API returns as a string.
type Amount string

// Float64 parses the amount as a floating point number. An empty amount is zero.
func (a Amount) Float64() (float64, error) {
	if a == "" {
		return 0, nil
	}

	return strconv.ParseFloat(string(a), 64)
}

// Int parses the amount as an integer. An empty amount is zero.
func (a Amount) Int() (int, error) {
	if a == "" {
		return 0, nil
	}

	return strconv.Atoi(string(a))
}

// String returns the raw value as returned by the API.
func (a Amount) String() string {
	return string(a)
}
