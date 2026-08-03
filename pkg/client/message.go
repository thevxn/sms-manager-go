package client

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"strings"
	"time"
)

const (
	// MaxRecipients is the maximum number of recipients of a single message.
	MaxRecipients = 10

	// MaxBatchSize is the maximum number of messages sent by a single SendBatch call.
	MaxBatchSize = 10

	// MaxBodyLength is the maximum length of a message body, longer texts get truncated by the API.
	MaxBodyLength = 1000
)

// Gateway is the SMS route used to deliver a message.
type Gateway string

const (
	// GatewayHigh is the default, high quality routing.
	GatewayHigh Gateway = "high"

	// GatewayLowcost is the low cost routing.
	GatewayLowcost Gateway = "lowcost"

	// GatewayDirect is used with a dedicated virtual number.
	GatewayDirect Gateway = "direct"

	// GatewayCustom is used with SIM hosting.
	GatewayCustom Gateway = "custom"

	// GatewaySimhost is used with SIM hosting.
	GatewaySimhost Gateway = "simhost"

	// GatewayGSM is used with an own GSM gateway.
	GatewayGSM Gateway = "gsm"
)

// Message tags. Multiple tags can be passed comma-separated, any other value
// groups messages together (e.g. by campaign ID).
const (
	// TagPromotional is the default tag, such messages are checked against the opt-out list.
	TagPromotional = "promotional"

	// TagTransactional marks transactional messages.
	TagTransactional = "transactional"

	// TagPriority marks priority messages.
	TagPriority = "priority"
)

// TTLCondition is evaluated once a flow step TTL expires. If the message has
// already reached the given state, no further flow step is used.
type TTLCondition string

const (
	TTLConditionSent      TTLCondition = "sent"
	TTLConditionDelivered TTLCondition = "delivered"
	TTLConditionSeen      TTLCondition = "seen"
)

// SMSType controls the unicode handling of an SMS flow step.
type SMSType string

const (
	// SMSTypeUTF keeps unicode characters in the message (the default).
	SMSTypeUTF SMSType = "utf"

	// SMSTypeSMS converts or strips unicode characters to keep the maximum SMS capacity.
	SMSTypeSMS SMSType = "sms"
)

// Weekday is a day of the week a delivery window applies to.
type Weekday string

const (
	Monday    Weekday = "monday"
	Tuesday   Weekday = "tuesday"
	Wednesday Weekday = "wednesday"
	Thursday  Weekday = "thursday"
	Friday    Weekday = "friday"
	Saturday  Weekday = "saturday"
	Sunday    Weekday = "sunday"
)

// Time marshals as an RFC 3339 timestamp in UTC, which is the format the API
// expects for scheduling. Use NewTime to build one.
type Time struct {
	time.Time
}

// NewTime wraps t for use as a scheduling timestamp.
func NewTime(t time.Time) *Time {
	return &Time{Time: t}
}

// MarshalJSON implements the json.Marshaler interface.
func (t Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.UTC().Format(time.RFC3339))
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *Time) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if raw == "" {
		return nil
	}

	parsed, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return err
	}

	t.Time = parsed
	return nil
}

// Recipient is a single message receiver.
type Recipient struct {
	// PhoneNumber is the recipient's number, ideally in the international E.164
	// format without the leading + or 00. Other formats are normalized by the API.
	PhoneNumber string `json:"phone_number"`
}

// Recipients converts phone numbers to a Recipient slice usable as Message.To.
func Recipients(phoneNumbers ...string) []Recipient {
	recipients := make([]Recipient, 0, len(phoneNumbers))

	for _, phoneNumber := range phoneNumbers {
		recipients = append(recipients, Recipient{PhoneNumber: phoneNumber})
	}

	return recipients
}

// Message is a single sending request of the JSON API v2.
type Message struct {
	// Body is the message text (max 1000 characters). It is compulsory either
	// here, or in a flow step --- a flow step definition takes precedence.
	Body string `json:"body,omitempty"`

	// To is a list of 1 to 10 recipients.
	To []Recipient `json:"to"`

	// Sender is the default sender for all flow steps --- an alphanumeric sender
	// name (max 11 characters, may require prior registration), or a dedicated
	// virtual number (without the leading +).
	Sender string `json:"sender,omitempty"`

	// Callback is the webhook target for delivery notifications. Pass a string
	// with the URL, or a struct/map with an extended callback configuration.
	Callback interface{} `json:"callback,omitempty"`

	// Tag groups messages together (e.g. by campaign ID). Special values are
	// TagPromotional (the default, checked against the opt-out list),
	// TagTransactional and TagPriority.
	Tag string `json:"tag,omitempty"`

	// Params carries special message parameters, e.g. automatic link shortening.
	Params *Params `json:"params,omitempty"`

	// DateTime schedules the message. Times more than a minute in the future are
	// scheduled, times in the past are sent immediately.
	DateTime *Time `json:"datetime,omitempty"`

	// DeliveryTime narrows the delivery down to an allowed time window.
	DeliveryTime *DeliveryTime `json:"delivery_time,omitempty"`

	// Flow is an ordered list of channels to attempt. Defaults to a single SMS
	// step on the GatewayHigh gateway when left empty.
	Flow []FlowStep `json:"flow,omitempty"`

	// Payload is arbitrary metadata returned back in delivery notifications and responses.
	Payload map[string]interface{} `json:"payload,omitempty"`
}

// Params carries special message parameters.
type Params struct {
	// URL enables automatic link shortening of links found in the message body.
	URL *URLShortening `json:"_url,omitempty"`
}

// URLShortening configures the automatic replacement of links in a message body.
type URLShortening struct {
	// Type is "full" to shorten complete URLs including the protocol only (the
	// default), or "all" to shorten every domain and address found.
	Type string `json:"type,omitempty"`

	// Domain is a custom domain used for the shortened links (has to be preconfigured).
	Domain string `json:"domain,omitempty"`

	// GetEnhance set to "id" appends the unique message_id to the shortened link
	// as a query parameter, so a click can be tied back to a recipient.
	GetEnhance string `json:"get_enhance,omitempty"`

	// GetName is the name of the query parameter used for GetEnhance (defaults to "_").
	GetName string `json:"get_name,omitempty"`

	// ExpiresAt is the UNIX timestamp of the shortened link expiration (defaults to 7 days).
	ExpiresAt int64 `json:"expires_at,omitempty"`
}

// DeliveryTime is the window a message may be delivered in. Sending outside of
// the window is postponed to the nearest valid time.
type DeliveryTime struct {
	// Days are the weekdays delivery is allowed on.
	Days []Weekday `json:"days,omitempty"`

	// Start is the beginning of the delivery window, e.g. "08:00".
	Start string `json:"start,omitempty"`

	// End is the end of the delivery window, e.g. "21:00".
	End string `json:"end,omitempty"`

	// TZ is the time zone of the delivery window, e.g. "Europe/Prague" (defaults to UTC).
	TZ string `json:"tz,omitempty"`
}

// FlowStep is a single channel attempt. Exactly one field has to be set --- use
// several steps in Message.Flow to define fallbacks between channels.
type FlowStep struct {
	SMS              *SMSStep              `json:"sms,omitempty"`
	Viber            *ViberStep            `json:"viber,omitempty"`
	WhatsAppText     *WhatsAppTextStep     `json:"whatsapp_text,omitempty"`
	WhatsAppTemplate *WhatsAppTemplateStep `json:"whatsapp_template,omitempty"`
	RCS              *RCSStep              `json:"rcs,omitempty"`
}

// SMSStep is an SMS channel attempt.
type SMSStep struct {
	// Body takes precedence over the Message.Body definition.
	Body string `json:"body,omitempty"`

	// Sender overrides Message.Sender for this step.
	Sender string `json:"sender,omitempty"`

	// Gateway is the SMS route used (defaults to GatewayHigh).
	Gateway Gateway `json:"gateway,omitempty"`

	// TTL is the time to live in minutes (min 0.5). Once it expires and
	// TTLCondition is not met, the next flow step is used.
	TTL float64 `json:"ttl,omitempty"`

	// TTLCondition is the state checked once the TTL expires (defaults to TTLConditionSent).
	TTLCondition TTLCondition `json:"ttl_condition,omitempty"`

	// Type controls the unicode handling (defaults to SMSTypeUTF).
	Type SMSType `json:"type,omitempty"`
}

// ViberStep is a Viber channel attempt.
type ViberStep struct {
	// Body takes precedence over the Message.Body definition.
	Body string `json:"body,omitempty"`

	// Sender has to be registered for Viber Business Messages beforehand.
	Sender string `json:"sender,omitempty"`

	// Buttons is an optional single call-to-action button.
	Buttons []ViberButton `json:"buttons,omitempty"`

	// TTL is the time to live in minutes (min 0.5).
	TTL float64 `json:"ttl,omitempty"`

	// TTLCondition is the state checked once the TTL expires (defaults to TTLConditionDelivered).
	TTLCondition TTLCondition `json:"ttl_condition,omitempty"`
}

// ViberButton is a call-to-action button of a Viber message.
type ViberButton struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

// WhatsAppTextStep is a free-form WhatsApp message. It can only be used when the
// recipient has already replied to a template message, or messaged first.
type WhatsAppTextStep struct {
	// Body takes precedence over the Message.Body definition.
	Body string `json:"body,omitempty"`

	// Sender is the Phone Number ID of a registered WhatsApp number.
	Sender string `json:"sender,omitempty"`
}

// WhatsAppTemplateStep sends a preapproved WhatsApp template.
type WhatsAppTemplateStep struct {
	// TemplateName is the name of a registered and approved template.
	TemplateName string `json:"template_name,omitempty"`

	// Sender is the Phone Number ID of a registered WhatsApp number.
	Sender string `json:"sender,omitempty"`

	// Params are the values of the template body parameters.
	Params []string `json:"params,omitempty"`

	// ParamsHeader are the values of the template header parameters.
	ParamsHeader []string `json:"params_header,omitempty"`

	// ParamsButtons are the values of the template button parameters.
	ParamsButtons []string `json:"params_buttons,omitempty"`

	// Language is the template language, e.g. "en" or "cs".
	Language string `json:"language,omitempty"`

	// TTL is the time to live in minutes (min 0.5).
	TTL float64 `json:"ttl,omitempty"`

	// TTLCondition is the state checked once the TTL expires (defaults to TTLConditionDelivered).
	TTLCondition TTLCondition `json:"ttl_condition,omitempty"`
}

// RCSStep is an RCS channel attempt.
type RCSStep struct {
	// Body takes precedence over the Message.Body definition.
	Body string `json:"body,omitempty"`

	// Sender is an RCS agent registered for the account.
	Sender string `json:"sender,omitempty"`

	// TTL is the time to live in minutes (min 0.5).
	TTL float64 `json:"ttl,omitempty"`

	// TTLCondition is the state checked once the TTL expires (defaults to TTLConditionDelivered).
	TTLCondition TTLCondition `json:"ttl_condition,omitempty"`
}

// SendResponse is the result of a sending request. A request can be accepted
// only partially, so both Accepted and Rejected have to be checked.
type SendResponse struct {
	// RequestID is the unique identifier of the whole request.
	RequestID string `json:"request_id"`

	// Accepted holds the accepted recipients (Send, SendPriority), or the
	// accepted messages of a batch (SendBatch).
	Accepted []AcceptedMessage `json:"accepted"`

	// Rejected holds the rejected recipients/messages.
	Rejected []RejectedMessage `json:"rejected"`
}

// AcceptedMessage is a single accepted recipient or batch item.
type AcceptedMessage struct {
	// Key is the index of the recipient in the Message.To array (Send,
	// SendPriority), or the index of the message in the batch (SendBatch).
	Key string `json:"key"`

	// MessageID is the unique identifier of the accepted message. With SendBatch
	// it identifies the whole request --- append "-<recipient index>" to get the
	// message ID of a single recipient.
	MessageID string `json:"message_id"`
}

// RejectedMessage is a single rejected recipient or batch item.
type RejectedMessage struct {
	// Key is the index of the rejected recipient or batch item.
	Key string `json:"key"`
}

// SimpleMessage is a single SMS sent through the simplified endpoint. Use
// Message with Send for the full feature set (flow, multiple recipients,
// delivery windows, ...).
type SimpleMessage struct {
	// PhoneNumber is the recipient's number.
	PhoneNumber string

	// Message is the message text.
	Message string

	// Sender is the sender name or virtual number (may require prior registration).
	Sender string

	// Gateway is the SMS route used (defaults to GatewayHigh).
	Gateway Gateway

	// DateTime schedules the message, in UTC.
	DateTime *time.Time

	// CustomID is an own message identifier. It is stored in the message payload
	// (as customID) and returned in delivery notifications.
	CustomID string
}

// SimpleSendResponse is the result of a SendSimple call.
type SimpleSendResponse struct {
	// RequestID is the unique identifier of the request.
	RequestID string `json:"request_id"`

	// MessageID is the unique identifier of the message. Empty when rejected.
	MessageID string `json:"message_id"`

	// Accepted is true when the message got accepted for sending.
	Accepted bool `json:"accepted"`
}

// Send delivers a message to up to 10 recipients.
//
// POST /message
func (c *Client) Send(ctx context.Context, msg Message) (*SendResponse, error) {
	return c.send(ctx, "/message", msg)
}

// SendPriority is Send using the priority queue. Use it for time critical
// messages such as one-time passwords or transactional notifications.
//
// POST /message/priority
func (c *Client) SendPriority(ctx context.Context, msg Message) (*SendResponse, error) {
	return c.send(ctx, "/message/priority", msg)
}

func (c *Client) send(ctx context.Context, path string, msg Message) (*SendResponse, error) {
	if err := msg.validate(); err != nil {
		return nil, err
	}

	response := &SendResponse{}
	if err := c.postJSON(ctx, c.jsonBaseURL, path, msg, response); err != nil {
		return nil, err
	}

	return response, nil
}

// SendBatch delivers up to 10 messages (up to 10 recipients each) in a single
// call. Unlike Send, the returned message IDs identify whole messages --- append
// "-<recipient index>" to address a single recipient.
//
// POST /messages
func (c *Client) SendBatch(ctx context.Context, msgs []Message) (*SendResponse, error) {
	if len(msgs) == 0 {
		return nil, errors.New("no messages to send")
	}
	if len(msgs) > MaxBatchSize {
		return nil, errors.New("too many messages in a batch, 10 at most")
	}

	for _, msg := range msgs {
		if err := msg.validate(); err != nil {
			return nil, err
		}
	}

	response := &SendResponse{}
	if err := c.postJSON(ctx, c.jsonBaseURL, "/messages", msgs, response); err != nil {
		return nil, err
	}

	return response, nil
}

// SendSimple delivers a single SMS through the simplified form-encoded
// endpoint. Note the message can be rejected while the call itself succeeds ---
// check SimpleSendResponse.Accepted.
//
// POST /simple/message
func (c *Client) SendSimple(ctx context.Context, msg SimpleMessage) (*SimpleSendResponse, error) {
	if msg.PhoneNumber == "" {
		return nil, errors.New("no recipient phone number given")
	}
	if msg.Message == "" {
		return nil, errors.New("no message text given")
	}
	if msg.Sender == "" {
		return nil, errors.New("no sender given")
	}

	params := url.Values{}
	params.Set("phone_number", msg.PhoneNumber)
	params.Set("message", msg.Message)
	params.Set("sender", msg.Sender)

	if msg.Gateway != "" {
		params.Set("gateway", string(msg.Gateway))
	}
	if msg.DateTime != nil {
		params.Set("datetime", msg.DateTime.UTC().Format(time.RFC3339))
	}
	if msg.CustomID != "" {
		params.Set("customID", msg.CustomID)
	}

	response := &SimpleSendResponse{}
	if err := c.postForm(ctx, c.jsonBaseURL, "/simple/message", params, response); err != nil {
		return nil, err
	}

	return response, nil
}

// validate catches the request errors the API would reject anyway.
func (m Message) validate() error {
	if len(m.To) == 0 {
		return errors.New("no recipients given")
	}
	if len(m.To) > MaxRecipients {
		return errors.New("too many recipients, 10 at most")
	}

	for _, recipient := range m.To {
		if strings.TrimSpace(recipient.PhoneNumber) == "" {
			return errors.New("empty recipient phone number")
		}
	}

	if m.Body != "" {
		return nil
	}

	// with no message body in the root object, at least one flow step has to carry one
	for _, step := range m.Flow {
		switch {
		case step.SMS != nil && step.SMS.Body != "":
			return nil
		case step.Viber != nil && step.Viber.Body != "":
			return nil
		case step.WhatsAppText != nil && step.WhatsAppText.Body != "":
			return nil
		case step.WhatsAppTemplate != nil && step.WhatsAppTemplate.TemplateName != "":
			return nil
		case step.RCS != nil && step.RCS.Body != "":
			return nil
		}
	}

	return errors.New("no message body given, set it either on the message or on a flow step")
}
