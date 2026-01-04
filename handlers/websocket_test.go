package handlers

import (
	"gotimedate/config"
	"testing"
)

func TestNewWSHandler(t *testing.T) {
	cfg := &config.Config{
		AllowedOrigins: []string{"http://localhost:3000", "http://localhost:8080"},
	}
	handler := NewWSHandler(cfg)

	if handler.timeService == nil {
		t.Error("Expected timeService to be initialized")
	}

	if handler.cfg == nil {
		t.Error("Expected config to be set")
	}
}

func TestWebSocketCORS(t *testing.T) {
	cfg := &config.Config{
		AllowedOrigins: []string{"http://localhost:3000"},
	}
	cfg.CompileOrigins()

	tests := []struct {
		origin   string
		expected bool
	}{
		{"", true},
		{"http://localhost:3000", true},
		{"http://localhost:8080", false},
		{"https://evil.com", false},
	}

	for _, test := range tests {
		result := cfg.IsOriginAllowed(test.origin)
		if result != test.expected {
			t.Errorf("IsOriginAllowed(%s) = %v, expected %v", test.origin, result, test.expected)
		}
	}
}
