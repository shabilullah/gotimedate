package config

import (
	"testing"
)

func TestIsOriginAllowed(t *testing.T) {
	cfg := &Config{
		AllowedOrigins: []string{
			"http://localhost:3000",
			"https://*.domain.com",
			"http://localhost:*",
		},
	}

	cfg.CompileOrigins()

	tests := []struct {
		name     string
		origin   string
		expected bool
	}{
		{"Empty origin", "", true},
		{"Exact match", "http://localhost:3000", true},
		{"Subdomain match", "https://app.domain.com", true},
		{"Deep subdomain match", "https://v1.api.domain.com", true},
		{"Subdomain mismatch", "https://domain.com.evil.com", false},
		{"Port wildcard match", "http://localhost:8080", true},
		{"Port wildcard match 2", "http://localhost:5173", true},
		{"Port wildcard mismatch host", "http://otherhost:8080", false},
		{"Not in list", "https://google.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cfg.IsOriginAllowed(tt.origin); got != tt.expected {
				t.Errorf("IsOriginAllowed(%q) = %v, want %v", tt.origin, got, tt.expected)
			}
		})
	}
}
