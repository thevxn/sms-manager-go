//package	"go.savla.dev/sms-manager"
package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

var (
	APIEndpoint string	= "https://http-api.smsmanager.cz/Send"
	APIKey string		= ""
	phoneNumber string	= ""
	message string		= ""
	messageID int		= 225153515
)

type Request struct {
	apiKey string `params: "apiKey"`
}

func main() {
	baseURL, err := url.Parse(APIEndpoint)
	if err != nil {
		log.Fatal(err)
	}
	params := url.Values{}

	// compulsory params
	params.Add("apikey", APIKey)
	params.Add("number", phoneNumber)
	params.Add("message", message)
	params.Add("gateway", "high")

	// optional params
	//params.Add("sender", "KPU")
	//params.Add("customid", string(messageID))
	//params.Add("time", "2023-06-26T13:30:00")
	//params.Add("expiration", "2023-06-26T14:00:00")

	baseURL.RawQuery = params.Encode()

	resp, err := http.Get(baseURL.String())
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	log.Println("status: " + resp.Status)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
	}

	log.Print(string(body))
}
