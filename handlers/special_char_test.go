package handlers

import (
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestTimeHandler_SpecialCharacterValidation(t *testing.T) {
	app := fiber.New()
	h := NewTimeHandler()
	app.Get("/api/v1/time", h.GetCurrentTime)

	t.Run("Reject SQL injection attempts", func(t *testing.T) {
		sqlInjectionAttempts := []string{
			"UTC; DROP TABLE users; --",
			"UTC OR 1=1",
			"UTC\\x27; DROP TABLE users; --",
			"UTC\\x22; DROP TABLE users; --",
			"UTC UNION SELECT * FROM users",
		}

		for _, maliciousInput := range sqlInjectionAttempts {
			req, _ := http.NewRequest("GET", "/api/v1/time?timezone="+url.QueryEscape(maliciousInput), nil)
			resp, _ := app.Test(req)

			if resp.StatusCode != http.StatusBadRequest {
				t.Errorf("expected status 400 for SQL injection attempt %q, got %v", maliciousInput, resp.StatusCode)
			}
		}
	})

	t.Run("Reject XSS attempts", func(t *testing.T) {
		xssAttempts := []string{
			"<script>alert(\"xss\")</script>",
			"javascript:alert(\"xss\")",
			"<img src=x onerror=alert(\"xss\")>",
			"UTC\"><script>alert(\"xss\")</script>",
		}

		for _, maliciousInput := range xssAttempts {
			req, _ := http.NewRequest("GET", "/api/v1/time?timezone="+url.QueryEscape(maliciousInput), nil)
			resp, _ := app.Test(req)

			if resp.StatusCode != http.StatusBadRequest {
				t.Errorf("expected status 400 for XSS attempt %q, got %v", maliciousInput, resp.StatusCode)
			}
		}
	})

	t.Run("Reject path traversal attempts", func(t *testing.T) {
		pathTraversalAttempts := []string{
			"../../../etc/passwd",
			"..\\..\\..\\windows\\system32",
			"/etc/passwd",
			"C:\\Windows\\System32",
			"UTC/../../../etc/passwd",
		}

		for _, maliciousInput := range pathTraversalAttempts {
			req, _ := http.NewRequest("GET", "/api/v1/time?timezone="+url.QueryEscape(maliciousInput), nil)
			resp, _ := app.Test(req)

			if resp.StatusCode != http.StatusBadRequest {
				t.Errorf("expected status 400 for path traversal attempt %q, got %v", maliciousInput, resp.StatusCode)
			}
		}
	})

	t.Run("Reject control characters", func(t *testing.T) {
		controlCharacterAttempts := []string{
			"UTC\x00null",
			"UTC\x01start of heading",
			"UTC\x1bescape",
			"UTC\r\nnewline",
			"UTC\ttab",
			"UTC\x0fshift in",
		}

		for _, maliciousInput := range controlCharacterAttempts {
			req, _ := http.NewRequest("GET", "/api/v1/time?timezone="+url.QueryEscape(maliciousInput), nil)
			resp, _ := app.Test(req)

			if resp.StatusCode != http.StatusBadRequest {
				t.Errorf("expected status 400 for control character attempt %q, got %v", maliciousInput, resp.StatusCode)
			}
		}
	})

	t.Run("Handle unicode characters gracefully", func(t *testing.T) {
		unicodeInputs := []string{
			"UTC🌍",       // emoji
			"UTC中文",      // chinese
			"UTCالعربية", // arabic
			"UTC🕐",       // clock emoji
		}

		for _, unicodeInput := range unicodeInputs {
			req, _ := http.NewRequest("GET", "/api/v1/time?timezone="+url.QueryEscape(unicodeInput), nil)
			resp, _ := app.Test(req)

			if resp.StatusCode != http.StatusBadRequest {
				body, _ := io.ReadAll(resp.Body)
				t.Logf("Unicode input %q returned status %v with body: %s", unicodeInput, resp.StatusCode, string(body))
			}
		}
	})
}
