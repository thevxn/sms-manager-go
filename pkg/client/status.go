package client

import (
	"context"
	"errors"
	"net/url"
)

// MessageStatus is the full state of a single sent message.
type MessageStatus struct {
	// StatusCode is the sending status code, StatusSent means sent.
	StatusCode StatusCode `json:"status_code"`

	// DeliveryCode is the delivery status code, DeliveryDelivered means
	// delivered, DeliverySeen means seen.
	DeliveryCode StatusCode `json:"delivery_code"`

	// StatusInfo is a human-readable status description, optionally prefixed
	// with a code in brackets. Use ParseStatusInfo to split it.
	StatusInfo string `json:"status_info"`

	// Tag is the tag the message was sent with.
	Tag string `json:"tag"`

	// Body is the sent message text.
	Body string `json:"body"`

	// To is the recipient as submitted.
	To Recipient `json:"to"`

	// RequestID is the identifier of the sending request the message belongs to.
	RequestID string `json:"request_id"`

	TimeCreated   *Time `json:"ts_created"`
	TimeSent      *Time `json:"ts_sent"`
	TimeDelivered *Time `json:"ts_delivered"`
	TimeSeen      *Time `json:"ts_seen"`
	TimeClicked   *Time `json:"ts_clicked"`
	TimeUpdated   *Time `json:"ts_updated"`

	// Flow is the original flow of the message.
	Flow []FlowStep `json:"flow"`

	// FlowIndex is the index of the flow step that was used.
	FlowIndex string `json:"flow_index"`

	// FlowStatus holds the status codes of the individual flow steps, keyed by step index.
	FlowStatus map[string]FlowStatus `json:"flow_status"`

	// Price is the price of the message in the account currency.
	Price *float64 `json:"price"`
}

// FlowStatus is the state of a single flow step.
type FlowStatus struct {
	StatusCode   StatusCode `json:"status_code"`
	DeliveryCode StatusCode `json:"delivery_code"`
}

// ListMessagesParams narrows down a ListMessages call.
type ListMessagesParams struct {
	// Date is the compulsory date to list the messages for, in the YYYY-MM-DD format.
	Date string

	// SubID is an optional subaccount ID.
	SubID string

	// Tag filters the messages by tag.
	Tag string

	// PhoneNumber filters the messages by the recipient's phone number.
	PhoneNumber string

	// Cursor is the NextCursor of a previous response, used for paging.
	Cursor string
}

// MessageList is a single page of sent messages.
type MessageList struct {
	// Messages are the per-recipient overview rows.
	Messages []SentMessage `json:"messages"`

	// NextCursor is to be passed back as ListMessagesParams.Cursor to fetch the
	// next page. Nil when there are no more pages.
	NextCursor *string `json:"next_cursor"`
}

// SentMessage is an overview row of a sent message.
type SentMessage struct {
	MessageID string `json:"message_id"`

	// MessageParts is the number of parts the message was split (and charged) into.
	MessageParts int `json:"message_parts"`

	TimeSent      *Time `json:"ts_sent"`
	TimeDelivered *Time `json:"ts_delivered"`
	TimeUpdated   *Time `json:"ts_updated"`
	TimeClicked   *Time `json:"ts_clicked"`
	TimeSeen      *Time `json:"ts_seen"`
	TimeReplied   *Time `json:"ts_replied"`
	TimeOptOut    *Time `json:"ts_optout"`

	Sent         int `json:"sent"`
	Delivered    int `json:"delivered"`
	Seen         int `json:"seen"`
	Clicked      int `json:"clicked"`
	Replied      int `json:"replied"`
	OptOut       int `json:"optout"`
	NonDelivered int `json:"nondelivered"`

	Body string `json:"body"`

	// Flow is the channel the message was sent over, e.g. "sms" or "viber".
	Flow string `json:"flow"`

	Sender      string `json:"sender"`
	PhoneNumber string `json:"phone_number"`

	PriceCZK Amount `json:"price_czk"`
	PriceEUR Amount `json:"price_eur"`
}

// RequestStatsParams identifies the sending request to get the statistics of.
type RequestStatsParams struct {
	// Date is the compulsory date of the request, in the YYYY-MM-DD format.
	Date string

	// RequestID is the compulsory identifier returned by a sending call.
	RequestID string

	// SubID is an optional subaccount ID.
	SubID string
}

// RequestStats are the aggregated statistics of a single sending request.
type RequestStats struct {
	CountDelivered    Amount `json:"count_delivered"`
	CountNonDelivered Amount `json:"count_nondelivered"`
	CountSent         Amount `json:"count_sent"`
	CountNotSent      Amount `json:"count_notsent"`
	CountReplied      Amount `json:"count_replied"`
	CountOptOut       Amount `json:"count_optout"`

	// MessagesSent is the number of message parts sent.
	MessagesSent Amount `json:"msg_sent"`

	PriceCZK Amount `json:"price_czk"`
	PriceEUR Amount `json:"price_eur"`

	FirstTimeSent *Time `json:"first_ts_sent"`
	LastTimeSent  *Time `json:"last_ts_sent"`
}

// MessageStatus returns the state of a single sent message. It returns
// ErrNotFound if no such message was sent by the calling API key.
//
// GET /message
func (c *Client) MessageStatus(ctx context.Context, messageID string) (*MessageStatus, error) {
	if messageID == "" {
		return nil, errors.New("no message ID given")
	}

	params := url.Values{}
	params.Set("id", messageID)

	var response struct {
		Success bool           `json:"success"`
		Message *MessageStatus `json:"message"`
	}

	if err := c.get(ctx, c.restBaseURL, "/message", params, &response); err != nil {
		return nil, err
	}

	if response.Message == nil {
		return nil, ErrNotFound
	}

	return response.Message, nil
}

// ListMessages returns a page of messages sent on the given date. Pass the
// returned MessageList.NextCursor back as ListMessagesParams.Cursor to page
// through the rest.
//
// GET /messages
func (c *Client) ListMessages(ctx context.Context, listParams ListMessagesParams) (*MessageList, error) {
	if listParams.Date == "" {
		return nil, errors.New("no date given")
	}

	params := url.Values{}
	params.Set("date", listParams.Date)

	if listParams.SubID != "" {
		params.Set("subid", listParams.SubID)
	}
	if listParams.Tag != "" {
		params.Set("tag", listParams.Tag)
	}
	if listParams.PhoneNumber != "" {
		params.Set("phone_number", listParams.PhoneNumber)
	}
	if listParams.Cursor != "" {
		params.Set("cursor", listParams.Cursor)
	}

	response := &MessageList{}
	if err := c.get(ctx, c.restBaseURL, "/messages", params, response); err != nil {
		return nil, err
	}

	return response, nil
}

// RequestStats returns the aggregated statistics of a single sending request.
// It returns ErrNotFound if there are no statistics for the given date and ID.
//
// GET /request
func (c *Client) RequestStats(ctx context.Context, statsParams RequestStatsParams) (*RequestStats, error) {
	if statsParams.Date == "" {
		return nil, errors.New("no date given")
	}
	if statsParams.RequestID == "" {
		return nil, errors.New("no request ID given")
	}

	params := url.Values{}
	params.Set("date", statsParams.Date)
	params.Set("id", statsParams.RequestID)

	if statsParams.SubID != "" {
		params.Set("subid", statsParams.SubID)
	}

	var response struct {
		Success bool          `json:"success"`
		Request *RequestStats `json:"request"`
	}

	if err := c.get(ctx, c.restBaseURL, "/request", params, &response); err != nil {
		return nil, err
	}

	if response.Request == nil {
		return nil, ErrNotFound
	}

	return response.Request, nil
}
