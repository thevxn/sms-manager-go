# sms-manager-go

[![Go Reference](https://pkg.go.dev/badge/go.vxn.dev/sms-manager.svg)](https://pkg.go.dev/go.vxn.dev/sms-manager)
[![Go Report Card](https://goreportcard.com/badge/go.vxn.dev/sms-manager)](https://goreportcard.com/report/go.vxn.dev/sms-manager)

Simple Go implementation of the [SMS Manager JSON API v2](https://smsmanager.cz/docs/api-reference/json-v2/overview) (message sending) and of the [REST API v1](https://smsmanager.cz/docs/api-reference/rest/overview) (account, message status, pricing).

The former HTTP API v1 (`http-api.smsmanager.cz`) is deprecated by the provider and is no longer implemented here --- see [migration](#migration-from-the-http-api-v1) below.

## usage

See the [example](/example) folder for the full client usage.

### in code

```go
package main

import (
	"context"
	"log"

	"go.vxn.dev/sms-manager/pkg/client"
)

func main() {
	c, err := client.New("your-api-key")
	if err != nil {
		log.Fatal(err)
	}

	resp, err := c.Send(context.Background(), client.Message{
		Body:   "hello world!",
		To:     client.Recipients("420777123456"),
		Sender: "SmsManager",
		Tag:    client.TagTransactional,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println(resp.Accepted[0].MessageID)
}
```

The API key is issued in the account settings at [app.smsmanager.com](https://app.smsmanager.com/api-cloud) and is sent in the `x-api-key` header by the client.

### methods

JSON API v2 (`https://api.smsmngr.com/v2`):

| method | endpoint | purpose |
|---|---|---|
| `Send` | `POST /message` | send a message to up to 10 recipients |
| `SendPriority` | `POST /message/priority` | the same over the priority queue (OTPs, transactional messages) |
| `SendBatch` | `POST /messages` | up to 10 messages (up to 100 recipients) in a single call |
| `SendSimple` | `POST /simple/message` | single SMS, form-encoded, no flow support |

REST API v1 (`https://rest-api.smsmngr.com/v1`):

| method | endpoint | purpose |
|---|---|---|
| `MessageStatus` | `GET /message` | full state of a single sent message |
| `ListMessages` | `GET /messages` | messages sent on a date, cursor-paged |
| `RequestStats` | `GET /request` | aggregated statistics of a sending request |
| `APIKeyInfo` | `GET /apikey/me` | key scope, senders, account type |
| `Credit` | `GET /credit` | credit balance and price levels |
| `Prices` | `GET /credit/prices` | price list by channel and destination |

### flow

Other channels than SMS (Viber, WhatsApp, RCS) and fallbacks between them are defined by `Message.Flow` --- the steps are attempted in order, the next one is used once a step TTL expires without its `TTLCondition` being met.

```go
resp, err := c.Send(ctx, client.Message{
	To: client.Recipients("420777123456"),
	Flow: []client.FlowStep{
		{Viber: &client.ViberStep{Body: "your order is ready", Sender: "MyCompany", TTL: 10, TTLCondition: client.TTLConditionDelivered}},
		{SMS: &client.SMSStep{Body: "your order is ready", Sender: "MyCompany", Gateway: client.GatewayHigh}},
	},
})
```

### errors

Non-2xx responses are returned as `*client.APIError` (HTTP status, API message, request ID). Lookups with no match return `client.ErrNotFound`. Note a sending call can succeed while some recipients are rejected --- always check `SendResponse.Rejected`, and `SimpleSendResponse.Accepted` with `SendSimple`.

## migration from the HTTP API v1

The old pipe-delimited API is gone, so the client API changed with it:

| v1 | v2 |
|---|---|
| `DoRequest(req, "Send", resp)` | `c.Send(ctx, client.Message{...})` |
| `DoRequest(req, "RequestStatus", resp)` | `c.MessageStatus(ctx, messageID)` |
| `DoRequest(req, "RequestList", resp)` | `c.ListMessages(ctx, client.ListMessagesParams{Date: ...})` / `c.RequestStats(ctx, ...)` |
| `DoRequest(req, "GetUserInfo", resp)` | `c.APIKeyInfo(ctx)` + `c.Credit(ctx)` |
| `DoRequest(req, "GetPrice", resp)` | `c.Prices(ctx)` (destination price list, not a per-message quote) |
| `Request.CustomID` | `Message.Payload` (or `SimpleMessage.CustomID`) |
| `Request.ExpirationTime` | flow step `TTL` + `TTLCondition` |
| numeric `RequestID` | UUID `request_id` and per-recipient `message_id` |
| `SMSProcessCode` / `SMSDeliveryCode` | `client.StatusCode` (`200` sent, `201` delivered, `202` seen) |

The JSON API v2 has no per-message price quote endpoint --- the price of a sent message is reported by `MessageStatus`, in delivery webhooks, and by `ListMessages`.
