package gcm

import (
	"fmt"
	"testing"
)

var VALID_API_KEY = "API_KEY"

func TestCreateClient(t *testing.T) {
	client := CreateClient("api_key", GATEWAY)

	if client.ApiKey != "api_key" {
		t.Errorf("Expect 'api_key' but found %s", client.ApiKey)
	}
	if client.HttpClient == nil {
		t.Errorf("Expect not nil but found nil")
	}
}

func TestSendWithInvalidGateway(t *testing.T) {
	message := CreateMessage([]string{"1", "2"})
	message.DryRun = true

	client := CreateClient("api_key", "http://localhost:8080")
	response := client.send(message)

	if response == nil {
		t.Errorf("Expect not nil but found nil")
	} else {
		if len(response.Results[0].Error) == 0 {
			t.Errorf("Expect not nil but found nil")
		}

		if response.Results[0].RegistrationId != message.RegistrationIds[0] {
			t.Errorf("Expect %s but found %s", message.RegistrationIds[0], response.Results[0].RegistrationId)
		}
	}
}

func TestSendWithInvalidApiKey(t *testing.T) {
	message := CreateMessage([]string{"1", "2"})
	message.DryRun = true

	client := CreateClient("api_key", GATEWAY)
	response := client.send(message)

	if response == nil {
		t.Errorf("Expect not nil but found nil")
	} else {
		if !(response.Results[0].Error == AUTHENTICATION_ERROR || response.Results[0].Error == TIMEOUT) {
			t.Errorf("Expect %s or %s but found %s", AUTHENTICATION_ERROR, TIMEOUT, response.Results[0].Error)
		}
	}
}

func TestSendWithInvalidRegistrationId(t *testing.T) {
	message := CreateMessage([]string{"1"})
	message.DryRun = true

	client := CreateClient(VALID_API_KEY, GATEWAY)
	response := client.send(message)

	if response == nil {
		t.Errorf("Expect not nil but found nil")
	} else {
		if response.Results[0].RegistrationId != message.RegistrationIds[0] {
			t.Errorf("Expect %s but found %s", message.RegistrationIds[0], response.Results[0].RegistrationId)
		}

		if response.Results[0].Error != INVALID_REGISTRATION_TOKEN {
			t.Errorf("Expect %s but found %s", INVALID_REGISTRATION_TOKEN, response.Results[0].Error)
		}
	}
}

func TestSend1000(t *testing.T) {
	registrationIds := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		registrationIds[i] = fmt.Sprintf("%d %d %d %d %d %d %d %d %d %d", i, i, i, i, i, i, i, i, i, i)
	}
	originalMessage := CreateMessage(registrationIds)
	originalMessage.DryRun = true

	client := CreateClient(VALID_API_KEY, GATEWAY)
	response := client.send(originalMessage)

	if len(response.Results) != len(registrationIds) {
		t.Errorf("Expect %d but found %d", len(registrationIds), len(response.Results))
	}
}

func TestSendMessage(t *testing.T) {
	registrationIds := make([]string, 1001)
	for i := 0; i < 1001; i++ {
		registrationIds[i] = fmt.Sprintf("%d %d %d %d %d %d %d %d %d %d", i, i, i, i, i, i, i, i, i, i)
	}
	originalMessage := CreateMessage(registrationIds)
	originalMessage.DryRun = true

	client := CreateClient(VALID_API_KEY, GATEWAY)
	responses := client.SendMessage(originalMessage)

	if len(responses) != 2 {
		t.Errorf("Expect 2 but found %d", len(responses))
	}

	if len(responses[0].Results) != 1000 {
		t.Errorf("Expect 1000 but found %d", len(responses[0].Results))
	}

	if len(responses[1].Results) != 1 {
		t.Errorf("Expect 1 but found %d", len(responses[1].Results))
	}
}
