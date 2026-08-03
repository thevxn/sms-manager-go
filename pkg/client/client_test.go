package client

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const testAPIKey = "test-key"

// newTestClient returns a Client with both API base URLs pointed at handler.
func newTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	c, err := New(testAPIKey, WithJSONAPIBaseURL(server.URL), WithRESTAPIBaseURL(server.URL))
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	return c
}

func TestNewRequiresAPIKey(t *testing.T) {
	if _, err := New(""); !errors.Is(err, ErrMissingAPIKey) {
		t.Errorf("expected ErrMissingAPIKey, got %v", err)
	}
}

func TestSend(t *testing.T) {
	var (
		gotPath   string
		gotKey    string
		gotType   string
		gotBody   map[string]interface{}
		sendTime  = time.Date(2025, 1, 11, 11, 0, 0, 0, time.FixedZone("CET", 3600))
		wantSent  = "2025-01-11T10:00:00Z"
		responses = `{"request_id":"66666666","accepted":[{"key":"0","message_id":"pppppppp"}],"rejected":[{"key":"1"}]}`
	)

	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotKey = r.Header.Get("x-api-key")
		gotType = r.Header.Get("Content-Type")

		body, _ := io.ReadAll(r.Body)
		if err := json.Unmarshal(body, &gotBody); err != nil {
			t.Errorf("request body is not JSON: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, responses)
	})

	resp, err := c.Send(context.Background(), Message{
		Body:     "hello world!",
		To:       Recipients("420777123456", "420777654321"),
		Sender:   "SmsManager",
		Tag:      TagTransactional,
		DateTime: NewTime(sendTime),
		Flow: []FlowStep{
			{SMS: &SMSStep{Gateway: GatewayHigh, TTL: 10, TTLCondition: TTLConditionSent}},
		},
		Payload: map[string]interface{}{"order_id": "ORD-1001"},
	})
	if err != nil {
		t.Fatalf("Send: %v", err)
	}

	if gotPath != "/message" {
		t.Errorf("path = %q, want /message", gotPath)
	}
	if gotKey != testAPIKey {
		t.Errorf("x-api-key = %q, want %q", gotKey, testAPIKey)
	}
	if gotType != "application/json" {
		t.Errorf("content type = %q, want application/json", gotType)
	}

	// the scheduled time has to be sent as an UTC RFC 3339 timestamp
	if got := gotBody["datetime"]; got != wantSent {
		t.Errorf("datetime = %v, want %v", got, wantSent)
	}

	recipients, ok := gotBody["to"].([]interface{})
	if !ok || len(recipients) != 2 {
		t.Fatalf("to = %v, want 2 recipients", gotBody["to"])
	}
	if got := recipients[0].(map[string]interface{})["phone_number"]; got != "420777123456" {
		t.Errorf("first recipient = %v, want 420777123456", got)
	}

	flow, ok := gotBody["flow"].([]interface{})
	if !ok || len(flow) != 1 {
		t.Fatalf("flow = %v, want a single step", gotBody["flow"])
	}
	if _, ok := flow[0].(map[string]interface{})["sms"]; !ok {
		t.Errorf("flow step = %v, want an sms step", flow[0])
	}

	if resp.RequestID != "66666666" {
		t.Errorf("request ID = %q, want 66666666", resp.RequestID)
	}
	if len(resp.Accepted) != 1 || resp.Accepted[0].MessageID != "pppppppp" {
		t.Errorf("accepted = %+v, want a single message ID pppppppp", resp.Accepted)
	}
	if len(resp.Rejected) != 1 || resp.Rejected[0].Key != "1" {
		t.Errorf("rejected = %+v, want the key 1", resp.Rejected)
	}
}

func TestSendPriorityPath(t *testing.T) {
	var gotPath string

	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		io.WriteString(w, `{"request_id":"1","accepted":[],"rejected":[]}`)
	})

	if _, err := c.SendPriority(context.Background(), Message{Body: "otp 123456", To: Recipients("420777123456")}); err != nil {
		t.Fatalf("SendPriority: %v", err)
	}

	if gotPath != "/message/priority" {
		t.Errorf("path = %q, want /message/priority", gotPath)
	}
}

func TestSendValidation(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Error("no request expected, the message should not pass validation")
	})

	tests := map[string]Message{
		"no recipients": {Body: "hello"},
		"no body":       {To: Recipients("420777123456")},
		"empty number":  {Body: "hello", To: []Recipient{{PhoneNumber: " "}}},
		"too many to":   {Body: "hello", To: Recipients("1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11")},
		"bodyless flow": {To: Recipients("420777123456"), Flow: []FlowStep{{SMS: &SMSStep{Sender: "SmsManager"}}}},
	}

	for name, msg := range tests {
		if _, err := c.Send(context.Background(), msg); err == nil {
			t.Errorf("%s: expected an error, got none", name)
		}
	}

	// a body on a flow step is enough
	if err := (Message{To: Recipients("420777123456"), Flow: []FlowStep{{SMS: &SMSStep{Body: "hello"}}}}).validate(); err != nil {
		t.Errorf("flow step body should be accepted: %v", err)
	}
}

func TestSendSimple(t *testing.T) {
	var gotForm string

	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Content-Type"); got != "application/x-www-form-urlencoded" {
			t.Errorf("content type = %q, want application/x-www-form-urlencoded", got)
		}

		body, _ := io.ReadAll(r.Body)
		gotForm = string(body)

		io.WriteString(w, `{"request_id":"11111111","message_id":"aaaaaaaa","accepted":true}`)
	})

	resp, err := c.SendSimple(context.Background(), SimpleMessage{
		PhoneNumber: "420777123456",
		Message:     "test message",
		Sender:      "MyCompany",
		CustomID:    "ORD-1001",
	})
	if err != nil {
		t.Fatalf("SendSimple: %v", err)
	}

	want := "customID=ORD-1001&message=test+message&phone_number=420777123456&sender=MyCompany"
	if gotForm != want {
		t.Errorf("form = %q, want %q", gotForm, want)
	}

	if !resp.Accepted || resp.MessageID != "aaaaaaaa" {
		t.Errorf("response = %+v, want an accepted message aaaaaaaa", resp)
	}
}

func TestAPIError(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"Message":"User is not authorized to access this resource with an explicit deny"}`)
	})

	_, err := c.Send(context.Background(), Message{Body: "hello", To: Recipients("420777123456")})

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected an APIError, got %v", err)
	}

	if apiErr.StatusCode != http.StatusBadRequest {
		t.Errorf("status code = %d, want 400", apiErr.StatusCode)
	}
	if apiErr.Message == "" {
		t.Error("expected the API error message to be parsed")
	}
}

func TestMessageStatus(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("id"); got != "e27ff0ac" {
			t.Errorf("id = %q, want e27ff0ac", got)
		}

		io.WriteString(w, `{"success":true,"message":{"status_code":"200","delivery_code":"201",
			"status_info":"[0] Delivered","tag":"newsletter","body":"Hello world",
			"to":{"phone_number":"420777123456"},"request_id":"bc36f3d1",
			"ts_sent":"2026-07-30T10:00:01Z","ts_delivered":"2026-07-30T10:00:04Z","ts_seen":null,
			"flow":[{"sms":{}}],"flow_index":"0","flow_status":{"0":{"status_code":"200","delivery_code":"201"}},
			"price":1.2}}`)
	})

	status, err := c.MessageStatus(context.Background(), "e27ff0ac")
	if err != nil {
		t.Fatalf("MessageStatus: %v", err)
	}

	if status.StatusCode != StatusSent || status.StatusCode.Text() != "sent" {
		t.Errorf("status code = %q (%s), want 200 (sent)", status.StatusCode, status.StatusCode.Text())
	}
	if status.DeliveryCode != DeliveryDelivered {
		t.Errorf("delivery code = %q, want 201", status.DeliveryCode)
	}
	if status.To.PhoneNumber != "420777123456" {
		t.Errorf("phone number = %q, want 420777123456", status.To.PhoneNumber)
	}
	if status.TimeDelivered == nil || status.TimeDelivered.Format(time.RFC3339) != "2026-07-30T10:00:04Z" {
		t.Errorf("delivery time = %v, want 2026-07-30T10:00:04Z", status.TimeDelivered)
	}
	if status.TimeSeen != nil {
		t.Errorf("seen time = %v, want nil", status.TimeSeen)
	}
	if status.Price == nil || *status.Price != 1.2 {
		t.Errorf("price = %v, want 1.2", status.Price)
	}
	if step, ok := status.FlowStatus["0"]; !ok || step.DeliveryCode != DeliveryDelivered {
		t.Errorf("flow status = %+v, want step 0 delivered", status.FlowStatus)
	}
}

func TestMessageStatusNotFound(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"success":true,"message":null}`)
	})

	if _, err := c.MessageStatus(context.Background(), "unknown"); !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestListMessagesPaging(t *testing.T) {
	var gotQueries []string

	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotQueries = append(gotQueries, r.URL.RawQuery)

		if r.URL.Query().Get("cursor") == "" {
			io.WriteString(w, `{"success":true,"messages":[{"message_id":"e27ff0ac","message_parts":1,
				"sent":1,"delivered":1,"body":"Hello world","flow":"sms","sender":"MyCompany",
				"phone_number":"420777123456","price_czk":"1.2","price_eur":"0.05"}],"next_cursor":"page2"}`)
			return
		}

		io.WriteString(w, `{"success":true,"messages":[],"next_cursor":null}`)
	})

	params := ListMessagesParams{Date: "2026-07-30", Tag: "newsletter"}

	first, err := c.ListMessages(context.Background(), params)
	if err != nil {
		t.Fatalf("ListMessages: %v", err)
	}

	if len(first.Messages) != 1 {
		t.Fatalf("messages = %d, want 1", len(first.Messages))
	}

	price, err := first.Messages[0].PriceCZK.Float64()
	if err != nil || price != 1.2 {
		t.Errorf("price = %v (%v), want 1.2", price, err)
	}

	if first.NextCursor == nil {
		t.Fatal("expected a next cursor")
	}

	params.Cursor = *first.NextCursor

	second, err := c.ListMessages(context.Background(), params)
	if err != nil {
		t.Fatalf("ListMessages (page 2): %v", err)
	}
	if second.NextCursor != nil {
		t.Errorf("next cursor = %v, want nil", *second.NextCursor)
	}

	if gotQueries[0] != "date=2026-07-30&tag=newsletter" {
		t.Errorf("first query = %q", gotQueries[0])
	}
	if gotQueries[1] != "cursor=page2&date=2026-07-30&tag=newsletter" {
		t.Errorf("second query = %q", gotQueries[1])
	}
}

func TestCreditBalance(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/credit" {
			t.Errorf("path = %q, want /credit", r.URL.Path)
		}

		io.WriteString(w, `{"success":true,"credit":{"priceLevels":{"sms_230":2},"credit_czk":"150.5",
			"credit_eur":"0","billing_id":"1234","vat":"1","currency":"czk","tz":"Europe/Prague"}}`)
	})

	credit, err := c.Credit(context.Background())
	if err != nil {
		t.Fatalf("Credit: %v", err)
	}

	balance, err := credit.Balance()
	if err != nil {
		t.Fatalf("Balance: %v", err)
	}
	if balance != 150.5 {
		t.Errorf("balance = %v, want 150.5", balance)
	}

	if credit.PriceLevels["sms_230"] != 2 {
		t.Errorf("price levels = %v, want sms_230 = 2", credit.PriceLevels)
	}
}

func TestAPIKeyInfo(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"success":true,"me":{"note":"eshop","account":"prepaid","scope":"full",
			"is_blocked":0,"senders":[{"flow":"sms","value":"MyCompany","name":"MyCompany"}],
			"currency":"czk","tz":"Europe/Prague"}}`)
	})

	info, err := c.APIKeyInfo(context.Background())
	if err != nil {
		t.Fatalf("APIKeyInfo: %v", err)
	}

	if info.Account != AccountPrepaid || info.Scope != ScopeFull {
		t.Errorf("account/scope = %q/%q, want prepaid/full", info.Account, info.Scope)
	}
	if len(info.Senders) != 1 || info.Senders[0].Value != "MyCompany" {
		t.Errorf("senders = %+v, want a single MyCompany sender", info.Senders)
	}
}

func TestPrices(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"success":true,"prices":{"sms":{"420":{"mcc":230,"price_czk":0.89,
			"price_eur":0.04,"operators":{"1":{"price_czk":0.89,"price_eur":0.04}}}}}}`)
	})

	prices, err := c.Prices(context.Background())
	if err != nil {
		t.Fatalf("Prices: %v", err)
	}

	price, ok := prices["sms"]["420"]
	if !ok {
		t.Fatal("no price for sms/420")
	}
	if price.MCC != 230 || price.PriceCZK != 0.89 {
		t.Errorf("price = %+v, want MCC 230 at 0.89 CZK", price)
	}
	if price.Operators["1"].PriceEUR != 0.04 {
		t.Errorf("operator price = %+v, want 0.04 EUR", price.Operators["1"])
	}
}

func TestParseStatusInfo(t *testing.T) {
	tests := []struct {
		info            string
		wantCode        string
		wantDescription string
	}{
		{"[307] Insufficient credit", "307", "Insufficient credit"},
		{"[368]", "368", ""},
		{"Unauthorized error", "", "Unauthorized error"},
		{"", "", ""},
	}

	for _, test := range tests {
		code, description := ParseStatusInfo(test.info)
		if code != test.wantCode || description != test.wantDescription {
			t.Errorf("ParseStatusInfo(%q) = %q/%q, want %q/%q", test.info, code, description, test.wantCode, test.wantDescription)
		}
	}

	if ResultCodeText["307"] != "insufficient credit" {
		t.Errorf("result code 307 = %q", ResultCodeText["307"])
	}
}
