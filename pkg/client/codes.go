package client

// SMSProcessCode indicates the state of SMS processing.
// https://smsmanager.cz/api/codes/#sms
type SMSProcessCode int

var ProcessCodeText []string = []string{
	"waiting for send",
	"being sent",
	"sent",
	"partly sent",
	"invalid request data",
	"sending failed",
	"cancelled",
	"rejected by system",
	"",
	"being processed",
}

// SMSDeliveryCode indicates the state of SMS delivery.
// https://smsmanager.cz/api/codes/#delivery
type SMSDeliveryCode int

var DeliveryCodeText []string = []string{
	"no delivery message",
	"being delivered",
	"delivered",
	"cannot be delivered",
	"unknown delivery state",
	"waiting for delivery",
	"delivery expired",
	"invalid phone number",
	"cannot be delivered (internal server errror)",
}

// ErrorCode indicates custom HTTP code error response by server (invalid params etc).
// https://smsmanager.cz/api/codes/#api
type ErrorCode int

var ErrorCodeText map[interface{}]string = map[interface{}]string{
	101:   "nonexistent data",
	102:   "wrong request data format",
	103:   "invalid username or password",
	104:   "invalid 'gateway' param",
	105:   "low credit",
	109:   "missing data",
	201:   "invalid phone number(s)",
	202:   "missing text, or text too long",
	203:   "invalid 'sender' param",
	"9xx": "internal server error",
}
