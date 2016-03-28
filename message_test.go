package gcm

import (
	"fmt"
	"reflect"
	"testing"
)

func TestCreateMessage(t *testing.T) {
	data := []struct {
		registrationIDs []string
		deviceIDs       []string

		result *Message
	}{
		{nil, nil, nil},
		{make([]string, 0), make([]string, 0), nil},
		// {[]string{"0", "1"}, nil, &Message{RegistrationIDs: []string{"0", "1"}, Priority: "high", DelayWhileIdle: false, Data: make(map[string]interface{})}},
	}

	for idx, test := range data {
		if m := CreateMessage(test.deviceIDs, test.registrationIDs); m != test.result {
			t.Errorf("[Test %d]:Expected nil but found not nil.", idx)
		}
	}

	// Test valid message
	message := CreateMessage(nil, []string{"0", "1"})
	if !reflect.DeepEqual(message.RegistrationIDs, []string{"0", "1"}) {
		t.Errorf("Expected %s but found %s.", []string{"0", "1"}, message.RegistrationIDs)
	}
	if message.DeviceIDs != nil {
		t.Error("Expected nil but found not nil.")
	}
	if message.CollapseKey != "" {
		t.Errorf("Expected nil but found %s.", message.CollapseKey)
	}
	if message.Priority != "high" {
		t.Errorf("Expected %s but found %s.", "high", message.CollapseKey)
	}
	if message.DryRun {
		t.Errorf("Expected %t but found %t.", false, message.DryRun)
	}
	if message.DelayWhileIdle {
		t.Errorf("Expected %t but found %t.", false, message.DelayWhileIdle)
	}
	if message.Data == nil || len(message.Data) != 0 {
		t.Error("Expected not nil but found nil.")
	}

	message = CreateMessage([]string{"0", "1"}, []string{"0", "1"})
	if !reflect.DeepEqual(message.DeviceIDs, []string{"0", "1"}) {
		t.Errorf("Expected %s but found %s.", []string{"0", "1"}, message.DeviceIDs)
	}
}

func TestSetField(t *testing.T) {
	message := CreateMessage(nil, []string{"0", "1"})

	// [Test 1] assign single key-value
	message.SetField("", "")
	if len(message.Data) != 0 {
		t.Errorf("Expected size is zero but found %d.", len(message.Data))
	}
	message.SetField("", "testValue")
	if len(message.Data) != 0 {
		t.Errorf("Expected size is zero but found %d.", len(message.Data))
	}

	message.SetField("testKey", "testValue")
	if len(message.Data) != 1 && message.Data["testKey"] != "testValue" {
		t.Errorf("Expected size is one but found %d.", len(message.Data))
	}

	// [Test 2] assign multiple keys-values
	data := []struct {
		key   string
		value interface{}
	}{
		{"", nil},
		{"key", nil},
		{"key", "value"},
		{"number", 100},
	}
	for _, test := range data {
		message.SetField(test.key, test.value)

		if message.Data[test.key] != test.value {
			t.Errorf("Expected %s but found %s", message.Data[test.key], test.value)
		}
	}
}

func TestEncode(t *testing.T) {
	registrationIDs := []string{"1"}
	message := CreateMessage(nil, registrationIDs)
	message.RestrictedPackage = "test"
	message.Data["key"] = "value"

	// [Test 1] encode with single registrationID
	messages := message.encode()
	if messages == nil {
		t.Error("Expected not nil but found nil")
	}
	if len(messages) != 1 {
		t.Errorf("Expected only 1 message after encoded but found %d", len(messages))
	}
	if len(messages) >= 1 {
		if messages[0].RestrictedPackage != message.RestrictedPackage {
			t.Errorf("Expected %s but found %s", message.RestrictedPackage, messages[0].RestrictedPackage)
		}
		if messages[0].Data["key"] != message.Data["key"] {
			t.Errorf("Expected %s but found %s", message.Data["key"], messages[0].Data["key"])
		}
	}

	// [Test 2] encode with 1000 registrationIDs
	registrationIDs = make([]string, 1000)
	for i := 0; i < 1000; i++ {
		registrationIDs[i] = fmt.Sprintf("%d", i)
	}
	message.RegistrationIDs = registrationIDs
	messages = message.encode()
	if len(messages) != 1 {
		t.Errorf("Expected only 1 message after encoded but found %d", len(messages))
	}

	// [Test 3] encode with 1001 registrationIDs
	registrationIDs = make([]string, 1001)
	deviceIDs := make([]string, 1001)
	for i := 0; i < 1001; i++ {
		deviceIDs[i] = fmt.Sprintf("%d", i)
		registrationIDs[i] = fmt.Sprintf("%d", i)
	}

	message.RegistrationIDs = registrationIDs
	message.DeviceIDs = deviceIDs
	messages = message.encode()
	if len(messages) != 2 {
		t.Errorf("Expected only 1 message after encoded but found %d", len(messages))
	}
	if messages[1].DeviceIDs[0] != "1000" {
		t.Errorf("Expected %s but found %s", "1000", messages[1].DeviceIDs[0])
	}
}
