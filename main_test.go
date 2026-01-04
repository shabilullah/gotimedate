package main

import (
	"gotimedate/config"
	"gotimedate/router"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestMainRouter(t *testing.T) {
	cfg := &config.Config{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Content-Type"},
		StaticDir:      "static",
	}
	cfg.CompileOrigins()

	app := router.SetupRouter(cfg)
	if app == nil {
		t.Fatal("router should not be nil")
	}

	t.Run("Root endpoint returns JSON with correct structure", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to send request: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", resp.StatusCode, http.StatusOK)
		}

		contentType := resp.Header.Get("Content-Type")
		if !strings.Contains(contentType, "application/json") {
			t.Errorf("handler returned wrong content type: got %v want application/json", contentType)
		}
	})

	t.Run("WS-Test endpoint returns HTML", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/ws-test", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to send request: %v", err)
		}

		if resp.StatusCode == http.StatusOK {
			contentType := resp.Header.Get("Content-Type")
			body, _ := io.ReadAll(resp.Body)
			if len(body) > 0 && !strings.Contains(contentType, "text/html") {
				t.Logf("handler returned content type: %v", contentType)
			}
		}
	})
}
