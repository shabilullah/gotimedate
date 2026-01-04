package services

import (
	"gotimedate/models"
	"testing"
	"time"
)

func TestTimeService_GetCurrentTime(t *testing.T) {
	s := NewTimeService()

	t.Run("Valid Timezone", func(t *testing.T) {
		tz := "Asia/Kuala_Lumpur"
		resp, err := s.GetCurrentTime(tz)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if resp.Timezone != tz {
			t.Errorf("expected timezone %s, got %s", tz, resp.Timezone)
		}
		if resp.Timestamp == "" {
			t.Error("expected non-empty timestamp")
		}
	})

	t.Run("Invalid Timezone", func(t *testing.T) {
		_, err := s.GetCurrentTime("Invalid/Zone")
		if err == nil {
			t.Error("expected error for invalid timezone, got nil")
		}
	})
}

func TestTimeService_ConvertTime(t *testing.T) {
	s := NewTimeService()

	t.Run("Valid Conversion", func(t *testing.T) {
		req := &models.TimeConvertRequest{
			Timestamp:    "2026-01-04T15:00:00Z",
			FromTimezone: "UTC",
			ToTimezone:   "Asia/Kuala_Lumpur",
		}
		resp, err := s.ConvertTime(req)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		// KL is UTC+8
		if resp.OffsetHours != 8 {
			t.Errorf("expected offset 8, got %f", resp.OffsetHours)
		}
	})
}

func TestTimeService_FormatTime(t *testing.T) {
	s := NewTimeService()
	now := time.Date(2026, 1, 4, 15, 4, 5, 0, time.UTC)

	t.Run("12hour format", func(t *testing.T) {
		got := s.FormatTime(now, "12hour")
		want := "3:04:05 PM"
		if got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})

	t.Run("24hour format", func(t *testing.T) {
		got := s.FormatTime(now, "24hour")
		want := "15:04:05"
		if got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})
}

func TestTimeService_GetAvailableTimezones(t *testing.T) {
	s := NewTimeService()
	zones := s.GetAvailableTimezones()
	if len(zones) == 0 {
		t.Error("expected at least one timezone")
	}

	foundKL := false
	for _, z := range zones {
		if z.Name == "Asia/Kuala_Lumpur" {
			foundKL = true
			break
		}
	}
	if !foundKL {
		t.Error("Asia/Kuala_Lumpur not found in available timezones")
	}
}
