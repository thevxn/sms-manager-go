package client

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"

	"go.savla.dev/sms-manager/pkg/config"
)

// composeURL is a helper function that parses params into a full url.URL type.
func composeURL(req config.Request, method string) (*url.URL, error) {
	if method == "" {
		return nil, errors.New("method string is empty")
	}

	baseURL, err := url.Parse(config.Endpoint + method)
	if err != nil {
		return nil, err
	}

	params := url.Values{}

	// compulsory params
	params.Add("apikey", req.APIKey)
	params.Add("number", req.PhoneNumber)
	params.Add("message", req.Message)
	params.Add("gateway", req.GatewayType)

	// optional params
	if req.RequestID != "" {
		params.Add("requestID", req.RequestID)
	}
	if req.Sender != "" {
		params.Add("sender", req.Sender)
	}
	if req.CustomID != "" {
		params.Add("customid", req.CustomID)
	}
	if req.SendTime != "" {
		params.Add("time", req.SendTime)
	}
	if req.ExpirationTime != "" {
		params.Add("expiration", req.ExpirationTime)
	}

	log.Println("method: " + method)

	// run the encoding process
	baseURL.RawQuery = params.Encode()

	return baseURL, nil
}

// getResponseBody does the actual HTTP GET request and returns read response body contents.
func getResponseBody(urlString string) ([]byte, error) {
	resp, err := http.Get(urlString)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	log.Println("status: " + resp.Status)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}
