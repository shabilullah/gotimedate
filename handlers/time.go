package handlers

import (
	"strings"
	"time"

	"gotimedate/models"
	"gotimedate/services"

	"github.com/gofiber/fiber/v2"
)

type TimeHandler struct {
	timeService *services.TimeService
}

func NewTimeHandler() *TimeHandler {
	return &TimeHandler{timeService: services.NewTimeService()}
}

// @Summary Get current time
// @Tags Time
// @Param timezone query string false "Timezone (default UTC)"
// @Success 200 {object} models.TimeResponse
// @Router /time [get]
func (h *TimeHandler) GetCurrentTime(c *fiber.Ctx) error {
	tz := c.Query("timezone", "UTC")
	resp, err := h.timeService.GetCurrentTime(tz)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.JSON(resp)
}

// @Summary Get time by timezone
// @Tags Time
// @Param timezone path string true "Timezone"
// @Success 200 {object} models.TimeResponse
// @Router /time/{timezone} [get]
func (h *TimeHandler) GetTimeByTimezone(c *fiber.Ctx) error {
	tz := c.Params("*")
	if tz == "" {
		return fiber.NewError(fiber.StatusBadRequest, "timezone is required")
	}
	tz = strings.TrimPrefix(tz, "/")
	resp, err := h.timeService.GetCurrentTime(tz)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.JSON(resp)
}

// @Summary Get available timezones
// @Tags Time
// @Success 200 {array} models.TimezoneInfo
// @Router /timezones [get]
func (h *TimeHandler) GetAvailableTimezones(c *fiber.Ctx) error {
	timezones := h.timeService.GetAvailableTimezones()
	return c.JSON(timezones)
}

// @Summary Convert time
// @Tags Time
// @Param request body models.TimeConvertRequest true "Conversion request"
// @Success 200 {object} models.TimeConvertResponse
// @Router /time/convert [post]
func (h *TimeHandler) ConvertTime(c *fiber.Ctx) error {
	var req models.TimeConvertRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}
	resp, err := h.timeService.ConvertTime(&req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.JSON(resp)
}

// @Summary Health check
// @Tags Health
// @Success 200 {object} models.HealthResponse
// @Router /health [get]
func (h *TimeHandler) HealthCheck(c *fiber.Ctx) error {
	return c.JSON(models.HealthResponse{Status: "healthy", Timestamp: time.Now(), Version: "1.0.0"})
}
