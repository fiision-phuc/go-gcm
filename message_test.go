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
	if !message.DelayWhileIdle {
		t.Errorf("Expected %t but found %t.", true, message.DelayWhileIdle)
	}
	if !(message.Data != nil && len(message.Data) != 0) {
		t.Error("Expected not nil but found nil.")
	}

	message = CreateMessage([]string{"0", "1"}, []string{"0", "1"})
	if !reflect.DeepEqual(message.DeviceIDs, []string{"0", "1"}) {
		t.Errorf("Expected %s but found %s.", []string{"0", "1"}, message.DeviceIDs)
	}
}

func TestSetField(t *testing.T) {
	message := CreateMessage(nil, []string{"0", "1"})

	message.SetField("", "")
	if !(message.Data != nil && len(message.Data) != 0) {
		t.Errorf("Expected size if zero but found %d.", len(message.Data))
	}
	message.SetField("testKey", "")
	if !(message.Data != nil && len(message.Data) != 0) {
		t.Errorf("Expected size if zero but found %d.", len(message.Data))
	}
	message.SetField("", "testValue")
	if !(message.Data != nil && len(message.Data) != 0) {
		t.Errorf("Expected size if zero but found %d.", len(message.Data))
	}

	message.SetField("testKey", "testValue")
	if !(message.Data != nil && len(message.Data) != 1 && message.Data["testKey"] != "testValue") {
		t.Errorf("Expected size if zero but found %d.", len(message.Data))
	}

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
			t.Errorf("Expect %s but found %s", message.Data[test.key], test.value)
		}
	}
}

func TestEncode(t *testing.T) {
	registrationIDs := []string{"1"}
	originalMessage := CreateMessage(registrationIDs, nil)
	originalMessage.RestrictedPackage = "test"
	originalMessage.Data["key"] = "value"

	messages := originalMessage.encode()
	if messages == nil {
		t.Error("Expect not nil but found nil")
	} else if len(messages) != 1 {
		t.Errorf("Expect only 1 message after encoded but found %d", len(messages))
	}

	if len(messages) >= 1 {
		if messages[0].RestrictedPackage != originalMessage.RestrictedPackage {
			t.Errorf("Expect %s but found %s", originalMessage.RestrictedPackage, messages[0].RestrictedPackage)
		}

		if messages[0].Data["key"] != originalMessage.Data["key"] {
			t.Errorf("Expect %s but found %s", originalMessage.Data["key"], messages[0].Data["key"])
		}
	}
}
func TestEncode1000(t *testing.T) {
	registrationIds := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		registrationIds[i] = fmt.Sprintf("%d", i)
	}
	originalMessage := CreateMessage(registrationIds, nil)

	messages := originalMessage.encode()
	if messages == nil {
		t.Error("Expect not nil but found nil")
	} else if len(messages) != 1 {
		t.Errorf("Expect only 1 message after encoded but found %d", len(messages))
	}
}
func TestEncode1001(t *testing.T) {
	registrationIds := make([]string, 1001)
	for i := 0; i < 1001; i++ {
		registrationIds[i] = fmt.Sprintf("%d", i)
	}
	originalMessage := CreateMessage(registrationIds, nil)

	messages := originalMessage.encode()
	if messages == nil {
		t.Error("Expect not nil but found nil")
	} else if len(messages) != 2 {
		t.Errorf("Expect only 1 message after encoded but found %d", len(messages))
	} else if len(originalMessage.RegistrationIDs) != 1001 {
		t.Error("Original registration Ids list had been edited")
	}
}
