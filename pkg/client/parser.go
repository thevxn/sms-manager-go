package client

import (
	"errors"
	"strconv"
	"strings"
)

func catchRequestError(parts []string, response *Response) (bool, error) {
	if parts[0] == "ERROR" {
		response.Message = parts[0]

		errorCode, err := strconv.Atoi(parts[1])
		if err != nil {
			return false, err
		}
		response.ErrorCode = ErrorCode(errorCode)

		response.ErrorMessage = parts[2]
		return false, nil
	}
	return true, nil
}

func parseSend(parts []string, response *Response) error {
	count := len(parts)

	if ok, err := catchRequestError(parts, response); !ok {
		return err
	}

	if count != 3 && count != 4 {
		return errors.New("invalid response field count")
	}

	response.Message = parts[0]

	requestID := parts[1]

	response.RequestID = requestID

	response.PhoneNumber = parts[2]

	if count == 4 {
		cleanPart := strings.Split(parts[3], "\n")
		customID, err := strconv.Atoi(cleanPart[0])
		if err != nil {
			return err
		}
		response.CustomID = customID
	}
	return nil
}

func parseRequestList(parts []string, response *Response) error {
	if ok, err := catchRequestError(parts, response); !ok {
		return err
	}

	// input "stream" needs to be rejoined, split by line break, and then parsed line by line
	joinedParts := strings.Join(parts, "|")
	splitParts := strings.Split(joinedParts, "\n")

	for _, line := range splitParts {
		newParts := strings.Split(line, "|")

		count := len(newParts)

		if count != 7 {
			// last run can be blank, so just skip it
			//return errors.New("invalid response field count")
			continue
		}

		requestID := newParts[0]

		response.RequestIDArray = append(response.RequestIDArray, requestID)

		response.GatewayTypeArray = append(response.GatewayTypeArray, newParts[1])
		response.TimeSentArray = append(response.TimeSentArray, newParts[2])
		response.TimeExpiryArray = append(response.TimeExpiryArray, newParts[3])
		response.SenderArray = append(response.SenderArray, newParts[4])

		remainCount, err := strconv.Atoi(newParts[5])
		if err != nil {
			//return err
			remainCount = 0
		}
		response.RemainCountArray = append(response.RemainCountArray, remainCount)

		processCode, err := strconv.Atoi(newParts[6])
		if err != nil {
			return err
		}
		response.ProcessCodeArray = append(response.ProcessCodeArray, SMSProcessCode(processCode))

	}
	return nil
}

func parseRequestStatus(parts []string, response *Response) error {
	count := len(parts)

	if ok, err := catchRequestError(parts, response); !ok {
		return err
	}

	if count != 4 {
		return errors.New("invalid response field count")
	}

	response.PhoneNumber = parts[0]

	processCode, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}
	response.ProcessCode = SMSProcessCode(processCode)

	deliveryCount, err := strconv.Atoi(parts[2])
	if err != nil {
		return err
	}
	response.DeliveryCount = deliveryCount

	deliveryCode, err := strconv.Atoi(parts[3])
	if err != nil {
		return err
	}
	response.DeliveryCode = SMSDeliveryCode(deliveryCode)

	return nil
}

func parseGetUserInfo(parts []string, response *Response) error {
	count := len(parts)

	if ok, err := catchRequestError(parts, response); !ok {
		return err
	}

	if count != 3 {
		return errors.New("invalid response field count")
	}

	// returns type float64!
	creditBalance, err := strconv.ParseFloat(parts[0], 32)
	if err != nil {
		return err
	}
	response.CreditBalance = creditBalance

	response.Sender = parts[1]
	response.GatewayType = parts[2]

	return nil
}

func parseGetPrice(parts []string, response *Response) error {
	count := len(parts)

	if ok, err := catchRequestError(parts, response); !ok {
		return err
	}

	if count != 5 {
		return errors.New("invalid response field count")
	}

	validReceiverCount, err := strconv.Atoi(parts[0])
	if err != nil {
		return err
	}
	response.ValidReceiverCount = validReceiverCount

	textSplitSMSCount, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}
	response.TextSplitSMSCount = textSplitSMSCount

	characterCount, err := strconv.Atoi(parts[2])
	if err != nil {
		return err
	}
	response.CharacterCount = characterCount

	// returns type float64!
	oneSMSPrice, err := strconv.ParseFloat(parts[3], 32)
	if err != nil {
		return err
	}
	response.OneSMSPrice = oneSMSPrice

	// wipe out newline character
	cleanPart := strings.Split(parts[4], "\n")
	// returns type float64!
	sumSMSPrice, err := strconv.ParseFloat(cleanPart[0], 32)
	if err != nil {
		return err
	}
	response.SumSMSPrice = sumSMSPrice

	return nil
}

// parseResponse distributes input raw <response> string, explodes it by a defined delimiter, and switches the operation to final subside parsing functions by <method> string.
func parseResponse(response string, method string, parsedResponse *Response) error {
	parts := strings.Split(response, "|")

	switch method {
	case "Send":
		if err := parseSend(parts, parsedResponse); err != nil {
			return err
		}
	case "RequestList":
		if err := parseRequestList(parts, parsedResponse); err != nil {
			return err
		}
	case "RequestStatus":
		if err := parseRequestStatus(parts, parsedResponse); err != nil {
			return err
		}
	case "GetUserInfo":
		if err := parseGetUserInfo(parts, parsedResponse); err != nil {
			return err
		}
	case "GetPrice":
		if err := parseGetPrice(parts, parsedResponse); err != nil {
			return err
		}
	}
	return nil
}
