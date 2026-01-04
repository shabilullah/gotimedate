package models

import "time"

type TimeResponse struct {
	Timestamp  string `json:"timestamp" example:"2024-01-03T14:30:45Z"`
	Timezone   string `json:"timezone" example:"UTC"`
	Unix       int64  `json:"unix" example:"1704315045"`
	UnixOffset int    `json:"unix_offset" example:"-18000"`
	Formatted  string `json:"formatted" example:"2:30:45 PM"`
	Date       string `json:"date" example:"Wednesday, January 3, 2024"`
}

type TimeConvertRequest struct {
	FromTimezone string `json:"from_timezone" example:"UTC"`
	ToTimezone   string `json:"to_timezone" example:"America/New_York"`
	Timestamp    string `json:"timestamp" example:"2024-01-03T14:30:45Z"`
}

type TimeConvertResponse struct {
	Original      TimeResponse `json:"original"`
	Converted     TimeResponse `json:"converted"`
	OffsetHours   float64      `json:"offset_hours" example:"-5.0"`
	OffsetMinutes int          `json:"offset_minutes" example:"-300"`
}

type TimezoneInfo struct {
	Name    string  `json:"name" example:"America/New_York"`
	Offset  float64 `json:"offset" example:"-5.0"`
	Current string  `json:"current" example:"2:30:45 PM"`
}

type WebSocketMessage struct {
	Type      string      `json:"type" example:"time_update"`
	Action    string      `json:"action,omitempty" example:"subscribe"`
	Timezone  string      `json:"timezone,omitempty" example:"America/New_York"`
	Format    string      `json:"format,omitempty" example:"12hour"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp string      `json:"timestamp,omitempty" example:"2024-01-03T14:30:45Z"`
}

type ErrorResponse struct {
	Error   string `json:"error" example:"Invalid timezone"`
	Message string `json:"message" example:"The specified timezone is not supported"`
	Code    int    `json:"code" example:"400"`
}

type HealthResponse struct {
	Status    string    `json:"status" example:"healthy"`
	Timestamp time.Time `json:"timestamp" example:"2024-01-03T14:30:45Z"`
	Version   string    `json:"version" example:"1.0.0"`
}

type TimeFormat struct {
	Name        string `json:"name" example:"ISO8601"`
	Description string `json:"description" example:"ISO 8601 format"`
	Example     string `json:"example" example:"2024-01-03T14:30:45Z"`
}
