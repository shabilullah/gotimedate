package services

import (
	"fmt"
	"gotimedate/models"
	"time"
)

type TimeService struct{}

func NewTimeService() *TimeService {
	return &TimeService{}
}

func (s *TimeService) GetCurrentTime(timezone string) (*models.TimeResponse, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, fmt.Errorf("invalid timezone: %s", timezone)
	}
	now := time.Now().In(loc)
	_, offset := now.Zone()
	return &models.TimeResponse{
		Timestamp:  now.Format(time.RFC3339),
		Timezone:   timezone,
		Unix:       now.Unix(),
		UnixOffset: offset,
		Formatted:  s.FormatTime(now, "12hour"),
		Date:       s.FormatDate(now),
	}, nil
}

func (s *TimeService) ConvertTime(req *models.TimeConvertRequest) (*models.TimeConvertResponse, error) {
	fromTime, err := time.Parse(time.RFC3339, req.Timestamp)
	if err != nil {
		return nil, fmt.Errorf("invalid timestamp format: %s", req.Timestamp)
	}
	fromLoc, err := time.LoadLocation(req.FromTimezone)
	if err != nil {
		return nil, fmt.Errorf("invalid from timezone: %s", req.FromTimezone)
	}
	toLoc, err := time.LoadLocation(req.ToTimezone)
	if err != nil {
		return nil, fmt.Errorf("invalid to timezone: %s", req.ToTimezone)
	}
	fromTimeInTZ := fromTime.In(fromLoc)
	toTimeInTZ := fromTime.In(toLoc)
	_, fromOffset := fromTimeInTZ.Zone()
	_, toOffset := toTimeInTZ.Zone()
	offsetSeconds := toOffset - fromOffset
	return &models.TimeConvertResponse{
		Original: models.TimeResponse{
			Timestamp:  fromTimeInTZ.Format(time.RFC3339),
			Timezone:   req.FromTimezone,
			Unix:       fromTimeInTZ.Unix(),
			UnixOffset: fromOffset,
			Formatted:  s.FormatTime(fromTimeInTZ, "12hour"),
			Date:       s.FormatDate(fromTimeInTZ),
		},
		Converted: models.TimeResponse{
			Timestamp:  toTimeInTZ.Format(time.RFC3339),
			Timezone:   req.ToTimezone,
			Unix:       toTimeInTZ.Unix(),
			UnixOffset: toOffset,
			Formatted:  s.FormatTime(toTimeInTZ, "12hour"),
			Date:       s.FormatDate(toTimeInTZ),
		},
		OffsetHours:   float64(offsetSeconds) / 3600.0,
		OffsetMinutes: offsetSeconds / 60,
	}, nil
}

func (s *TimeService) GetAvailableTimezones() []models.TimezoneInfo {
	zones := []string{
		"UTC",
		"Africa/Cairo", "Africa/Casablanca", "Africa/Johannesburg", "Africa/Lagos", "Africa/Nairobi",
		"America/Anchorage", "America/Argentina/Buenos_Aires", "America/Bogota", "America/Caracas",
		"America/Chicago", "America/Denver", "America/Halifax", "America/Los_Angeles",
		"America/Mexico_City", "America/New_York", "America/Phoenix", "America/Santiago", "America/Sao_Paulo",
		"Asia/Bangkok", "Asia/Dubai", "Asia/Hong_Kong", "Asia/Istanbul", "Asia/Jakarta",
		"Asia/Jerusalem", "Asia/Kabul", "Asia/Karachi", "Asia/Kolkata", "Asia/Kuala_Lumpur", "Asia/Manila",
		"Asia/Seoul", "Asia/Shanghai", "Asia/Singapore", "Asia/Taipei", "Asia/Tehran", "Asia/Tokyo",
		"Atlantic/Azores", "Atlantic/Cape_Verde",
		"Australia/Adelaide", "Australia/Brisbane", "Australia/Darwin", "Australia/Melbourne", "Australia/Perth", "Australia/Sydney",
		"Europe/Amsterdam", "Europe/Athens", "Europe/Berlin", "Europe/Brussels", "Europe/Budapest",
		"Europe/Dublin", "Europe/Lisbon", "Europe/London", "Europe/Luxembourg", "Europe/Madrid",
		"Europe/Moscow", "Europe/Oslo", "Europe/Paris", "Europe/Prague", "Europe/Rome",
		"Europe/Stockholm", "Europe/Vienna", "Europe/Warsaw", "Europe/Zurich",
		"Pacific/Auckland", "Pacific/Fiji", "Pacific/Guam", "Pacific/Honolulu", "Pacific/Pago_Pago",
	}

	var result []models.TimezoneInfo
	for _, tz := range zones {
		if loc, err := time.LoadLocation(tz); err == nil {
			now := time.Now().In(loc)
			_, offset := now.Zone()
			result = append(result, models.TimezoneInfo{
				Name:    tz,
				Offset:  float64(offset) / 3600.0,
				Current: s.FormatTime(now, "12hour"),
			})
		}
	}
	return result
}

func (s *TimeService) GetTimeFormats() []models.TimeFormat {
	return []models.TimeFormat{
		{Name: "ISO8601", Description: "ISO 8601 format", Example: time.RFC3339},
		{Name: "12hour", Description: "12-hour format", Example: "2:30:45 PM"},
		{Name: "24hour", Description: "24-hour format", Example: "14:30:45"},
	}
}

func (s *TimeService) FormatTime(t time.Time, format string) string {
	switch format {
	case "12hour":
		return t.Format("3:04:05 PM")
	case "24hour":
		return t.Format("15:04:05")
	default:
		return t.Format(time.RFC3339)
	}
}

func (s *TimeService) FormatDate(t time.Time) string {
	return t.Format("Monday, January 2, 2006")
}
