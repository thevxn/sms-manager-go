package config

// https://smsmanager.cz/api/http
type Request struct {
	// API key provided by SMS manager we b interface.
	APIKey string `param:"apikey"`

	// Phone number to send the SMS to.
	PhoneNumber string `param:"number"`

	// Message is the very content of such SMS.
	Message string `param:"message"`

	// GatewayType is the specification of gateway service/tariff used.
	GatewayType string `param:"gateway" default:"high"`

	// RequestID is to be parsed from sent SMS request string returned.
	RequestID string `param:"requestID" default:0`

	// Sender is the optional specification of a sender, telephone number or authorized text.
	Sender string `param:"sender"`

	// CustomID is a optional message index number.
	CustomID string `param:"customid"`

	// SendTime is the optional specification of the message sending time --- message scheduling; format 2023-12-31T23:59:59.
	SendTime string `param:"time"`

	// ExpirationTime is the optional specification of the message expiration time --- SMS won't be sent after this time; format 2023-12-31T23:59:59.
	ExpirationTime string `param:"expiration"`
}

var (
	Endpoint string = "https://http-api.smsmanager.cz/"
)
