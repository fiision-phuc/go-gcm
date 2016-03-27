package gcm

// Message describes a gcm message.
type Message struct {
	RegistrationIDs []string `json:"registration_ids"`
	DeviceIDs       []string `json:"_"`

	CollapseKey       string `json:"collapse_key,omitempty"`
	Priority          string `json:"priority,omitempty"`
	DryRun            bool   `json:"dry_run,omitempty"`
	DelayWhileIdle    bool   `json:"delay_while_idle,omitempty"`
	RestrictedPackage string `json:"restricted_package_name,omitempty"`

	Data map[string]interface{} `json:"data,omitempty"`
}

// CreateMessage returns a default message.
func CreateMessage(deviceIDs []string, registrationIDs []string) *Message {
	/* Condition validation */
	if len(registrationIDs) == 0 {
		return nil
	}

	return &Message{
		RegistrationIDs: registrationIDs,
		DeviceIDs:       deviceIDs,
		Priority:        "high",

		Data: make(map[string]interface{}),
	}
}

// SetField adds custom key-value pair to the message's payload.
func (m *Message) SetField(key string, value interface{}) {
	/* Condition validation */
	if len(key) == 0 || value == nil {
		return
	}
	m.Data[key] = value
}

// encode enforces the message. Split original message into multiple messages if required.
func (m *Message) encode() []*Message {
	length := len(m.RegistrationIDs)
	maxIDs := 1000

	/* Condition validation: return if number of devices is less than 1000 */
	if length <= maxIDs {
		return []*Message{m}
	}

	// Calculate step
	remain := length % maxIDs
	counter := (length - remain) / maxIDs
	if remain > 0 {
		counter++
	}

	// Create message collection
	messages := make([]*Message, counter)
	for i := 0; i < counter; i++ {
		strIdx := i * maxIDs
		endIdx := strIdx + maxIDs

		/* Condition validation: Validate upper bound */
		if endIdx > length {
			endIdx = length
		}
		message := *m
		message.RegistrationIDs = m.RegistrationIDs[strIdx:endIdx]
		if m.DeviceIDs != nil {
			message.DeviceIDs = m.DeviceIDs[strIdx:endIdx]
		}

		messages[i] = &message
	}
	return messages
}
