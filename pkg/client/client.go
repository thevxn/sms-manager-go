package client

import (
	"go.savla.dev/sms-manager/pkg/config"
)

func DoRequest(req config.Request, method string, response *Response) error {
	baseURL, err := composeURL(req, method)
	if err != nil {
		return err
	}

	body, err := getResponseBody(baseURL.String())
	if err != nil {
		return err
	}

	if err = parseResponse(string(body), method, response); err != nil {
		return err
	}

	return nil
}
