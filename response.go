package gcm

// Error returned from Google.
const (
	AuthenticationError       = "AuthenticationError"
	DeviceMessageRateExceeded = "DeviceMessageRateExceeded"
	InternalServerError       = "InternalServerError"
	InvalidDataKey            = "InvalidDataKey"
	InvalidJSON               = "InvalidJSON"
	InvalidPackageName        = "InvalidPackageName"
	InvalidRegistrationToken  = "InvalidRegistration"
	InvalidTimeToLive         = "InvalidTtl"
	MessageTooBig             = "MessageTooBig"
	MismatchSenderID          = "MismatchSenderId"
	MissingRegistrationToken  = "MissingRegistration"
	TopicsMessageRateExceeded = "TopicsMessageRateExceeded"
	Timeout                   = "Unavailable"
	UnregisteredDevice        = "NotRegistered"
)

// Response describes a response from Google.
type Response struct {
	MulticastID  int64 `json:"multicast_id"`
	Success      int   `json:"success"`
	Failure      int   `json:"failure"`
	CanonicalIDs int   `json:"canonical_ids"`

	Results []Result `json:"results"`
}

// Result describes a result for individual message.
type Result struct {
	DeviceID       string
	RegistrationID string

	MessageID         string `json:"message_id,omitempty"`
	Error             string `json:"error,omitempty"`
	NewRegistrationID string `json:"registration_id,omitempty"`
}

// ShouldPostpone validates if should not send a new message for a period of time.
func (r *Result) ShouldPostpone() bool {
	if r.Error == Timeout || r.Error == InternalServerError || r.Error == DeviceMessageRateExceeded || r.Error == TopicsMessageRateExceeded {
		return true
	}
	return false
}

// ShouldRemove validates error to decide if registrationId should be deleted or not.
func (r *Result) ShouldRemove() bool {
	if r.Error == UnregisteredDevice {
		return true
	}
	return false
}

// ShouldUpdate validates if registrationId should be update or not.
func (r *Result) ShouldUpdate() bool {
	if len(r.MessageID) > 0 && len(r.NewRegistrationID) > 0 {
		return true
	}
	return false
}
