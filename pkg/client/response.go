package client

type Response struct {
	// Message is human-readable response code. Equals "OK" or "ERROR".
	Message string

	// RequestID code can be used as reference to just sent SMS.
	RequestID int

	// RequestIDArray is an array of RequestID for RequestList method response parsing.
	RequestIDArray []int

	// PhoneNumber is a number for SMS delivery.
	PhoneNumber string

	// CustomID is a custom code set by client for SMS process reference.
	CustomID int

	// GatewayType is the name of used SMS gateway tariff.
	GatewayType string

	// GatewayTypeArray is an array of GatewayType/strings for RequestList method response parsing.
	GatewayTypeArray []string

	// TimeSentArray is an array of datetime strings for RequestList method response parsing.
	TimeSentArray []string

	// TimeExpiryArray is an array of datetime strings for RequestList method response parsing.
	TimeExpiryArray []string

	// Sender is a label of sender name set for given request (has to be approved, so mainly defaults to SmsManager or info-sms).
	Sender string

	// SenderArray is an array of string for RequestList method response parsing.
	SenderArray []string

	// RemainCountArray is an array of integers for RequestList method response parsing.
	// It tells the number of remaining recipients waiting to be served.
	RemainCountArray []int

	// DeliveryCount is a total sum of delivery acknowledgement (ACK) messages sent.
	DeliveryCount int

	// CreditBalance is a total amount of available prepaid financial resources.
	CreditBalance float64

	// ValidReceiverCount is a total number of valid phone numbers from given request (PhoneNumber request param).
	ValidReceiverCount int

	// TextSplitSMSCount is a total count of SMS messages theoretically sent from givem request (mainly driven by Message request param).
	TextSplitSMSCount int

	// CharacterCount indicates the sum of requested Message text length.
	CharacterCount int

	// OneSMSPrice is a price for one SMS theoretically sent from given input (mainly Message and PhoneNumber request params).
	OneSMSPrice float64

	// SumSMSPrice is a total price for all theoretically sent SMS messages from given input (mainly Message and PhoneNumber request params).
	SumSMSPrice float64

	// ProcessCode is a code returned by API with encoded meaning. Indicates the state of SMS processing.
	ProcessCode SMSProcessCode

	// ProcessCodeArray is an array of ProcessCode items for RequestList method parsing.
	ProcessCodeArray []SMSProcessCode

	// DeliveryCode is a code returned by API with encodedd meaning. Indicates the state of SMS delivery.
	DeliveryCode SMSDeliveryCode

	// ErrorCode is an encoded error type returned by some madmaniery happening between client and server.
	ErrorCode ErrorCode

	// ErrorMessage is a string containing the returned raw human-readable error message.
	ErrorMessage string
}
