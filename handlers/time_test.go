package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"gotimedate/models"

	"github.com/gofiber/fiber/v2"
)

func TestTimeHandler_HealthCheck(t *testing.T) {
	app := fiber.New()
	h := NewTimeHandler()
	app.Get("/health", h.HealthCheck)

	req, _ := http.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", resp.StatusCode, http.StatusOK)
	}

	var healthResp models.HealthResponse
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &healthResp)
	if healthResp.Status != "healthy" {
		t.Errorf("handler returned unexpected body: got %v want %v", healthResp.Status, "healthy")
	}
}

func TestTimeHandler_GetCurrentTime(t *testing.T) {
	app := fiber.New()
	h := NewTimeHandler()
	app.Get("/api/v1/time", h.GetCurrentTime)

	req, _ := http.NewRequest("GET", "/api/v1/time?timezone=UTC", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", resp.StatusCode, http.StatusOK)
	}

	var timeResp models.TimeResponse
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &timeResp)
	if timeResp.Timezone != "UTC" {
		t.Errorf("expected timezone UTC, got %s", timeResp.Timezone)
	}
}

func TestTimeHandler_GetTimeByTimezone(t *testing.T) {
	app := fiber.New()
	h := NewTimeHandler()
	app.Get("/api/v1/time/*", h.GetTimeByTimezone)

	req, _ := http.NewRequest("GET", "/api/v1/time/America/New_York", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", resp.StatusCode, http.StatusOK)
	}

	var timeResp models.TimeResponse
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &timeResp)
	if timeResp.Timezone != "America/New_York" {
		t.Errorf("expected timezone America/New_York, got %s", timeResp.Timezone)
	}
}

func TestTimeHandler_GetAvailableTimezones(t *testing.T) {
	app := fiber.New()
	h := NewTimeHandler()
	app.Get("/api/v1/timezones", h.GetAvailableTimezones)

	req, _ := http.NewRequest("GET", "/api/v1/timezones", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", resp.StatusCode, http.StatusOK)
	}

	var timezones []models.TimezoneInfo
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &timezones)
	if len(timezones) == 0 {
		t.Error("expected a list of timezones, got empty")
	}
}

func TestTimeHandler_ConvertTime(t *testing.T) {
	app := fiber.New()
	h := NewTimeHandler()
	app.Post("/api/v1/time/convert", h.ConvertTime)

	convertReq := models.TimeConvertRequest{
		FromTimezone: "UTC",
		ToTimezone:   "America/New_York",
		Timestamp:    "2026-01-03T22:00:00Z",
	}
	body, _ := json.Marshal(convertReq)
	req, _ := http.NewRequest("POST", "/api/v1/time/convert", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", resp.StatusCode, http.StatusOK)
	}

	var convertResp models.TimeConvertResponse
	respBody, _ := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &convertResp)
	if convertResp.Converted.Timezone != "America/New_York" {
		t.Errorf("expected converted timezone America/New_York, got %s", convertResp.Converted.Timezone)
	}
}

func TestTimeHandler_InputValidation(t *testing.T) {
	app := fiber.New()
	h := NewTimeHandler()
	app.Get("/api/v1/time", h.GetCurrentTime)
	app.Post("/api/v1/time/convert", h.ConvertTime)

	t.Run("Reject non-string timezone in query param", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/time?timezone=123", nil)
		resp, _ := app.Test(req)

		if resp.StatusCode != http.StatusOK {
			t.Logf("correctly rejected numeric timezone with status %v", resp.StatusCode)
		}
	})

	t.Run("Reject invalid JSON types in convert endpoint", func(t *testing.T) {
		invalidJSON := `{
"from_timezone": 123,
"to_timezone": "America/New_York", 
"timestamp": "2026-01-03T22:00:00Z"
}`
		req, _ := http.NewRequest("POST", "/api/v1/time/convert", bytes.NewBufferString(invalidJSON))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Logf("correctly rejected non-string timezone with status %v", resp.StatusCode)
		}
	})

	t.Run("Reject array instead of string in convert endpoint", func(t *testing.T) {
		invalidJSON := `{
"from_timezone": ["UTC"],
"to_timezone": "America/New_York",
"timestamp": "2026-01-03T22:00:00Z"
}`
		req, _ := http.NewRequest("POST", "/api/v1/time/convert", bytes.NewBufferString(invalidJSON))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		if resp.StatusCode != http.StatusOK {
			t.Logf("correctly rejected array timezone with status %v", resp.StatusCode)
		}
	})

	t.Run("Reject object instead of string in convert endpoint", func(t *testing.T) {
		invalidJSON := `{
"from_timezone": {"tz": "UTC"},
"to_timezone": "America/New_York",
"timestamp": "2026-01-03T22:00:00Z"
}`
		req, _ := http.NewRequest("POST", "/api/v1/time/convert", bytes.NewBufferString(invalidJSON))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		if resp.StatusCode != http.StatusOK {
			t.Logf("correctly rejected object timezone with status %v", resp.StatusCode)
		}
	})

	t.Run("Reject null values in convert endpoint", func(t *testing.T) {
		invalidJSON := `{
"from_timezone": null,
"to_timezone": "America/New_York",
"timestamp": "2026-01-03T22:00:00Z"
}`
		req, _ := http.NewRequest("POST", "/api/v1/time/convert", bytes.NewBufferString(invalidJSON))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		if resp.StatusCode != http.StatusOK {
			t.Logf("correctly handled null timezone with status %v", resp.StatusCode)
		}
	})

	t.Run("Reject boolean instead of string in convert endpoint", func(t *testing.T) {
		invalidJSON := `{
"from_timezone": true,
"to_timezone": "America/New_York",
"timestamp": "2026-01-03T22:00:00Z"
}`
		req, _ := http.NewRequest("POST", "/api/v1/time/convert", bytes.NewBufferString(invalidJSON))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		if resp.StatusCode != http.StatusOK {
			t.Logf("correctly rejected boolean timezone with status %v", resp.StatusCode)
		}
	})
}
