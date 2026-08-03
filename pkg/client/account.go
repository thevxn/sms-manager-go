package client

import (
	"context"
	"errors"
)

// Account types.
const (
	AccountPrepaid  = "prepaid"
	AccountPostpaid = "postpaid"
)

// API key permission scopes. ScopeFull and ScopeREST may change settings,
// ScopeAPI may only send messages.
const (
	ScopeFull = "full"
	ScopeREST = "rest-api"
	ScopeAPI  = "api"
)

// Account currencies.
const (
	CurrencyCZK = "czk"
	CurrencyEUR = "eur"
)

// APIKeyInfo is the profile of the calling API key.
type APIKeyInfo struct {
	// Note is the own note attached to the API key.
	Note string `json:"note"`

	// DefaultCallbackURL is the default webhook target for delivery notifications.
	DefaultCallbackURL string `json:"default_callback_url"`

	// Account is the account type, AccountPrepaid or AccountPostpaid.
	Account string `json:"account"`

	// Scope is the permission scope of the key, ScopeFull, ScopeREST or ScopeAPI.
	Scope string `json:"scope"`

	// ShareSenders is "1" when the key shares its senders with its subordinate keys.
	ShareSenders string `json:"share_senders"`

	// IsBlocked is 1 when the key is blocked.
	IsBlocked int `json:"is_blocked"`

	// Senders are the senders registered for the key.
	Senders []Sender `json:"senders"`

	// VAT is "1" when the account is a VAT payer.
	VAT string `json:"vat"`

	// Currency is the account currency, CurrencyCZK or CurrencyEUR.
	Currency string `json:"currency"`

	// TZ is the account time zone, e.g. "Europe/Prague".
	TZ string `json:"tz"`
}

// Sender is a sender registered for an API key.
type Sender struct {
	// Flow is the channel of the sender, e.g. "sms", "viber" or "whatsapp".
	Flow string `json:"flow"`

	// Value is the sender name or number to be used as Message.Sender.
	Value string `json:"value"`

	// Name is the display name of the sender.
	Name string `json:"name"`
}

// Credit is the credit state of the calling API key.
type Credit struct {
	// PriceLevels maps "<platform>_<mcc>" to the index of the reached discount
	// level (only entries with an index above 0 are present).
	PriceLevels map[string]int `json:"priceLevels"`

	// CreditCZK is the credit balance in CZK.
	CreditCZK Amount `json:"credit_czk"`

	// CreditEUR is the credit balance in EUR.
	CreditEUR Amount `json:"credit_eur"`

	// BillingID is the billing ID of the key.
	BillingID string `json:"billing_id"`

	// VAT is "1" when the account is a VAT payer.
	VAT string `json:"vat"`

	// Currency is the account currency, CurrencyCZK or CurrencyEUR.
	Currency string `json:"currency"`

	// TZ is the account time zone, e.g. "Europe/Prague".
	TZ string `json:"tz"`
}

// Balance returns the credit balance in the account currency.
func (c Credit) Balance() (float64, error) {
	if c.Currency == CurrencyEUR {
		return c.CreditEUR.Float64()
	}

	return c.CreditCZK.Float64()
}

// Prices is a price list keyed by channel ("sms", "smsdirect", "smscustom",
// "viber") and by country calling code.
type Prices map[string]map[string]DestinationPrice

// DestinationPrice is the price of a single channel and destination.
type DestinationPrice struct {
	// MCC is the official mobile country code of the destination.
	MCC int `json:"mcc"`

	PriceCZK float64 `json:"price_czk"`
	PriceEUR float64 `json:"price_eur"`

	// Operators holds per-operator prices keyed by the operator code, if any.
	Operators map[string]OperatorPrice `json:"operators"`
}

// OperatorPrice is the price of a single operator of a destination.
type OperatorPrice struct {
	PriceCZK float64 `json:"price_czk"`
	PriceEUR float64 `json:"price_eur"`
}

// APIKeyInfo returns the profile of the calling API key --- its senders, scope,
// account type and whether it is blocked.
//
// GET /apikey/me
func (c *Client) APIKeyInfo(ctx context.Context) (*APIKeyInfo, error) {
	var response struct {
		Success bool        `json:"success"`
		Me      *APIKeyInfo `json:"me"`
	}

	if err := c.get(ctx, c.restBaseURL, "/apikey/me", nil, &response); err != nil {
		return nil, err
	}

	if response.Me == nil {
		return nil, errors.New("no api key info returned")
	}

	return response.Me, nil
}

// Credit returns the credit balance and the reached price levels.
//
// GET /credit
func (c *Client) Credit(ctx context.Context) (*Credit, error) {
	var response struct {
		Success bool    `json:"success"`
		Credit  *Credit `json:"credit"`
	}

	if err := c.get(ctx, c.restBaseURL, "/credit", nil, &response); err != nil {
		return nil, err
	}

	if response.Credit == nil {
		return nil, errors.New("no credit info returned")
	}

	return response.Credit, nil
}

// Prices returns the price list by channel and destination.
//
// GET /credit/prices
func (c *Client) Prices(ctx context.Context) (Prices, error) {
	var response struct {
		Success bool   `json:"success"`
		Prices  Prices `json:"prices"`
	}

	if err := c.get(ctx, c.restBaseURL, "/credit/prices", nil, &response); err != nil {
		return nil, err
	}

	return response.Prices, nil
}
