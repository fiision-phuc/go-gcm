package gcm

import (
	"encoding/json"
	"testing"
)

func TestResponse(t *testing.T) {
	data := `{ "multicast_id": 216,
  			   "success": 3,
  			   "failure": 3,
  			   "canonical_ids": 1,
  			   "results": [
    			 { "message_id": "1:0408" },
    			 { "error": "Unavailable" },
    			 { "error": "InvalidRegistration" },
    			 { "message_id": "1:1516" },
    			 { "message_id": "1:2342", "registration_id": "32" },
    			 { "error": "NotRegistered"}
  			   ]
			 }`

	response := Response{}
	json.Unmarshal([]byte(data), &response)

	if response.MulticastID != 216 {
		t.Errorf("Expect %d but found %d", 216, response.MulticastID)
	}
	if response.Success != 3 {
		t.Errorf("Expect %d but found %d", 3, response.Success)
	}
	if response.Failure != 3 {
		t.Errorf("Expect %d but found %d", 3, response.Success)
	}
	if response.CanonicalIDs != 1 {
		t.Errorf("Expect %d but found %d", 1, response.CanonicalIDs)
	}

	if response.Results[0].ShouldPostpone() {
		t.Errorf("Expect %t but found %t", false, response.Results[0].ShouldPostpone())
	}
	if response.Results[0].ShouldUpdate() {
		t.Errorf("Expect %t but found %t", false, response.Results[0].ShouldUpdate())
	}
	if response.Results[0].ShouldRemove() {
		t.Errorf("Expect %t but found %t", false, response.Results[0].ShouldRemove())
	}

	if !response.Results[1].ShouldPostpone() {
		t.Errorf("Expect true but found %t", response.Results[1].ShouldPostpone())
	}
	if !response.Results[4].ShouldUpdate() {
		t.Errorf("Expect true but found %t", response.Results[4].ShouldUpdate())
	}
	if !response.Results[5].ShouldRemove() {
		t.Errorf("Expect true but found %t", response.Results[5].ShouldRemove())
	}
}
