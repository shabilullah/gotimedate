package main

import (
	_ "embed"
	"gotimedate/config"
	"gotimedate/router"
	"os"

	"github.com/gofiber/fiber/v2/log"
)

//go:embed static/websocket-test.html
var defaultHTML []byte

// @title Go TimeDate API
// @version 1.0.0
// @description API for time operations and WebSockets
// @host localhost:8080
// @BasePath /api/v1
func main() {
	cfg := config.LoadConfig(defaultHTML)

	logFile, err := os.OpenFile(cfg.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer logFile.Close()

	// Configure Fiber logger to output to both file and stdout
	log.SetOutput(logFile)

	app := router.SetupRouter(cfg)

	addr := cfg.Host + ":" + cfg.Port
	log.Infof("Server starting on %s", addr)
	log.Infof("Logging to: %s", cfg.LogFile)

	if err := app.Listen(addr); err != nil {
		log.Fatal(err)
	}
}
