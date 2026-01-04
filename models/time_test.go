package models

import (
	"encoding/json"
	"testing"
)

func TestTimeResponse_JSON(t *testing.T) {
	resp := TimeResponse{
		Timestamp: "2026-01-04T15:00:00Z",
		Timezone:  "UTC",
		Unix:      1767538800,
		Formatted: "3:00:00 PM",
		Date:      "Sunday, January 4, 2026",
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal TimeResponse: %v", err)
	}

	var unmarshaled TimeResponse
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal TimeResponse: %v", err)
	}

	if unmarshaled.Timestamp != resp.Timestamp {
		t.Errorf("expected timestamp %s, got %s", resp.Timestamp, unmarshaled.Timestamp)
	}
}

func TestTimeConvertRequest_JSON(t *testing.T) {
	req := TimeConvertRequest{
		FromTimezone: "UTC",
		ToTimezone:   "Asia/Kuala_Lumpur",
		Timestamp:    "2026-01-04T15:00:00Z",
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("failed to marshal TimeConvertRequest: %v", err)
	}

	var unmarshaled TimeConvertRequest
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal TimeConvertRequest: %v", err)
	}

	if unmarshaled.FromTimezone != req.FromTimezone {
		t.Errorf("expected from_timezone %s, got %s", req.FromTimezone, unmarshaled.FromTimezone)
	}
}
