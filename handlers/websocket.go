package handlers

import (
	"gotimedate/config"
	"gotimedate/models"
	"gotimedate/services"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/websocket/v2"
)

type WSHandler struct {
	timeService *services.TimeService
	cfg         *config.Config
}

func NewWSHandler(cfg *config.Config) *WSHandler {
	return &WSHandler{
		timeService: services.NewTimeService(),
		cfg:         cfg,
	}
}

func (h *WSHandler) ServeHTTP(c *websocket.Conn) {
	tz := "UTC"
	format := "12hour"
	stop := make(chan bool)

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				resp, _ := h.timeService.GetCurrentTime(tz)
				// Parse the timestamp from resp to get the timezone-converted time
				if parsedTime, err := time.Parse(time.RFC3339, resp.Timestamp); err == nil {
					resp.Formatted = h.timeService.FormatTime(parsedTime, format)
				}
				msg := models.WebSocketMessage{
					Type:      "time_update",
					Data:      resp,
					Timestamp: time.Now().Format(time.RFC3339),
				}
				if err := c.WriteJSON(msg); err != nil {
					log.Errorf("WebSocket write error: %v", err)
					return
				}
			case <-stop:
				return
			}
		}
	}()

	for {
		var msg models.WebSocketMessage
		if err := c.ReadJSON(&msg); err != nil {
			close(stop)
			break
		}
		if msg.Action == "subscribe" {
			if msg.Timezone != "" {
				tz = msg.Timezone
			}
			if msg.Format != "" {
				format = msg.Format
			}
		}
	}
}
