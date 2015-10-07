package gcm

//https://developers.google.com/cloud-messaging/http
//https://developers.google.com/cloud-messaging/http-server-ref#downstream-http-messages-json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	GATEWAY = "https://android.googleapis.com/gcm/send"

	MAX_BACKOFF_DELAY     = 1024000
	BACKOFF_INITIAL_DELAY = 1000
)

type Client struct {
	ApiKey     string
	Gateway    string
	HttpClient *http.Client
}

// MARK: Struct's constructors
func CreateClient(apiKey string, gateway string) *Client {
	return &Client{
		ApiKey:     apiKey,
		Gateway:    gateway,
		HttpClient: http.DefaultClient,
	}
}

// MARK: Struct's public functions
func (m *Client) SendMessage(templateMessage *Message) []*Response {
	/* Condition validation */
	if templateMessage == nil {
		return nil
	}

	// Encode template message
	messages := templateMessage.Encode()
	responses := make([]*Response, len(messages))

	for idx, message := range messages {
		responses[idx] = m.send(message)
	}
	return responses
}

// MARK: Struct's private functions
func (m *Client) send(message *Message) *Response {
	// Encode Json
	data, _ := json.Marshal(message)

	// Prepare request
	request, _ := http.NewRequest("POST", m.Gateway, bytes.NewBuffer(data))
	request.Header.Add("Authorization", fmt.Sprintf("key=%s", m.ApiKey))
	request.Header.Add("Content-Type", "application/json")

	// Send request
	httpResponse, _ := m.HttpClient.Do(request)
	if httpResponse == nil {
		response := &Response{
			MulticastId:  -1,
			Success:      0,
			Failure:      len(message.RegistrationIds),
			CanonicalIds: 0,
			Results:      make([]Result, len(message.RegistrationIds)),
		}

		for idx, registrationId := range message.RegistrationIds {
			response.Results[idx] = Result{Error: TIMEOUT, RegistrationId: registrationId}
		}

		return response
	}
	defer httpResponse.Body.Close()

	// Analyze response status
	if httpResponse.StatusCode != http.StatusOK {
		response := &Response{
			MulticastId:  -1,
			Success:      0,
			Failure:      len(message.RegistrationIds),
			CanonicalIds: 0,
			Results:      make([]Result, len(message.RegistrationIds)),
		}

		// Define error message
		errorMessage := ""
		if httpResponse.StatusCode == http.StatusUnauthorized {
			errorMessage = AUTHENTICATION_ERROR
		} else {
			errorMessage = INTERNAL_SERVER_ERROR
		}

		// Update result response
		for idx, registrationId := range message.RegistrationIds {
			response.Results[idx] = Result{Error: errorMessage, RegistrationId: registrationId}
		}
		return response
	} else {
		body, err := ioutil.ReadAll(httpResponse.Body)

		// Validate response data
		if err == nil {
			response := Response{}
			err = json.Unmarshal(body, &response)

			// Update result response
			if err == nil {
				for idx, registrationId := range message.RegistrationIds {
					response.Results[idx].RegistrationId = registrationId
				}
				return &response
			}
		}

		// Manual create response
		response := &Response{
			MulticastId:  -1,
			Success:      0,
			Failure:      len(message.RegistrationIds),
			CanonicalIds: 0,
			Results:      make([]Result, len(message.RegistrationIds)),
		}

		// Update result response
		for idx, registrationId := range message.RegistrationIds {
			response.Results[idx] = Result{Error: INVALID_JSON, RegistrationId: registrationId}
		}
		return response
	}
}
