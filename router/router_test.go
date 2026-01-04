package router

import (
	"gotimedate/config"
	"net/http"
	"testing"
)

func TestSetupRouter(t *testing.T) {
	cfg := &config.Config{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Content-Type"},
		StaticDir:      "static",
	}
	cfg.CompileOrigins()

	app := SetupRouter(cfg)
	if app == nil {
		t.Fatal("router should not be nil")
	}

	t.Run("Health Check Route", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/health", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status OK, got %v", resp.Status)
		}
	})

	t.Run("API Time Route", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/time", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status OK, got %v", resp.Status)
		}
	})

	t.Run("Invalid Route", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/invalid-route-123", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("expected status NotFound, got %v", resp.Status)
		}
	})
}
