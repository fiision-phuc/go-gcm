package gcm

import (
	"fmt"
	"testing"
)

// var apiKey = "API_KEY"
var apiKey = "AIzaSyBlsxfKNJCjGK1XSiI19K41X7dMlSeIuzU"

func TestCreateClient(t *testing.T) {
	client := CreateClient("api_key", Gateway)

	if client.APIKey != "api_key" {
		t.Errorf("Expected \"api_key\" but found \"%s\".", client.APIKey)
	}
	if client.HTTPClient == nil {
		t.Errorf("Expected not nil but found nil.")
	}
}

func TestSendMessage(t *testing.T) {
	message := CreateMessage([]string{"1", "2"}, []string{"3", "4"})
	message.DryRun = true

	// [Test 1] send with invalid api key
	client := CreateClient("api_key", Gateway)
	responses := client.SendMessage(message)

	if responses == nil || len(responses) == 0 {
		t.Errorf("Expected not nil but found nil.")
	} else {
		if len(responses[0].Results[0].Error) == 0 {
			t.Errorf("Expected not nil but found nil.")
		}
		if !(responses[0].Results[0].Error == AuthenticationError || responses[0].Results[0].Error == Timeout) {
			t.Errorf("Expected \"%s\" or \"%s\" but found \"%s\".", AuthenticationError, Timeout, responses[0].Results[0].Error)
		}

		if responses[0].Results[0].DeviceID != message.DeviceIDs[0] {
			t.Errorf("Expect \"%s\" but found \"%s\".", message.DeviceIDs[0], responses[0].Results[0].DeviceID)
		}
		if responses[0].Results[0].RegistrationID != message.RegistrationIDs[0] {
			t.Errorf("Expect \"%s\" but found \"%s\".", message.RegistrationIDs[0], responses[0].Results[0].RegistrationID)
		}
	}

	// [Test 2] send with invalid gateway
	client = CreateClient("api_key", "example.com")
	responses = client.SendMessage(message)

	if responses == nil || len(responses) == 0 {
		t.Errorf("Expected not nil but found nil.")
	} else {
		if responses[0].MulticastID != -1 {
			t.Errorf("Expected %d but found %d.", -1, responses[0].MulticastID)
		}
		if responses[0].Success != 0 {
			t.Errorf("Expected %d but found %d.", 0, responses[0].Success)
		}
		if responses[0].Failure != 2 {
			t.Errorf("Expected %d but found %d.", 2, responses[0].Failure)
		}

		if len(responses[0].Results[0].Error) == 0 {
			t.Errorf("Expected not nil but found nil.")
		}
		if responses[0].Results[0].Error != Timeout {
			t.Errorf("Expected \"%s\" or \"%s\" but found \"%s\".", AuthenticationError, Timeout, responses[0].Results[0].Error)
		}

		if responses[0].Results[0].DeviceID != message.DeviceIDs[0] {
			t.Errorf("Expect \"%s\" but found \"%s\".", message.DeviceIDs[0], responses[0].Results[0].DeviceID)
		}
		if responses[0].Results[0].RegistrationID != message.RegistrationIDs[0] {
			t.Errorf("Expect \"%s\" but found \"%s\".", message.RegistrationIDs[0], responses[0].Results[0].RegistrationID)
		}
	}

	// [Test 3] send with invalid registrationID
	client = CreateClient(apiKey, Gateway)
	responses = client.SendMessage(message)

	if responses == nil || len(responses) == 0 {
		t.Errorf("Expected not nil but found nil")
	} else {
		if responses[0].Results[0].DeviceID != message.DeviceIDs[0] {
			t.Errorf("Expect \"%s\" but found \"%s\".", message.DeviceIDs[0], responses[0].Results[0].DeviceID)
		}
		if responses[0].Results[0].RegistrationID != message.RegistrationIDs[0] {
			t.Errorf("Expected \"%s\" but found \"%s\".", message.RegistrationIDs[0], responses[0].Results[0].RegistrationID)
		}
		if responses[0].Results[0].Error != InvalidRegistrationToken {
			t.Errorf("Expected \"%s\" but found \"%s\".", InvalidRegistrationToken, responses[0].Results[0].Error)
		}
	}

	// [Test 4] send nil message
	client = CreateClient(apiKey, Gateway)
	responses = client.SendMessage(nil)

	if responses != nil || len(responses) != 0 {
		t.Errorf("Expected nil but found not nil.")
	}

	// [Test 5] send 1000 messages
	registrationIDs := make([]string, 1000)
	deviceIDs := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		registrationIDs[i] = fmt.Sprintf("%d %d %d %d %d %d %d %d %d %d", i, i, i, i, i, i, i, i, i, i)
		deviceIDs[i] = fmt.Sprintf("%d %d %d %d %d %d %d %d %d %d", i, i, i, i, i, i, i, i, i, i)
	}
	message.RegistrationIDs = registrationIDs
	message.DeviceIDs = deviceIDs
	message.DryRun = true

	client = CreateClient(apiKey, Gateway)
	responses = client.SendMessage(message)

	if len(responses) != 1 {
		t.Errorf("Expect %d but found %d", 1, len(responses))
	}
	if len(responses[0].Results) != 1000 {
		t.Errorf("Expect %d but found %d", 1000, len(responses[0].Results))
	}

	// [Test 6] send 1001 messages
	registrationIDs = make([]string, 1001)
	deviceIDs = make([]string, 1001)
	for i := 0; i < 1001; i++ {
		registrationIDs[i] = fmt.Sprintf("%d %d %d %d %d %d %d %d %d %d", i, i, i, i, i, i, i, i, i, i)
		deviceIDs[i] = fmt.Sprintf("%d %d %d %d %d %d %d %d %d %d", i, i, i, i, i, i, i, i, i, i)
	}
	message.RegistrationIDs = registrationIDs
	message.DeviceIDs = deviceIDs
	message.DryRun = true

	client = CreateClient(apiKey, Gateway)
	responses = client.SendMessage(message)

	if len(responses) != 2 {
		t.Errorf("Expect %d but found %d", 2, len(responses))
	}
	if len(responses[0].Results) != 1000 {
		t.Errorf("Expect %d but found %d", 1000, len(responses[0].Results))
	}
	if len(responses[1].Results) != 1 {
		t.Errorf("Expect %d but found %d", 1, len(responses[1].Results))
	}

	if responses[1].Results[0].DeviceID != message.DeviceIDs[1000] {
		t.Errorf("Expect \"%s\" but found \"%s\".", message.DeviceIDs[1000], responses[1].Results[0].DeviceID)
	}
	if responses[1].Results[0].RegistrationID != message.RegistrationIDs[1000] {
		t.Errorf("Expected \"%s\" but found \"%s\".", message.RegistrationIDs[1000], responses[1].Results[0].RegistrationID)
	}
}
