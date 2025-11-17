package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/neoway/golang-validation-documents/internal/config"
)

type StatusHandler struct{}

func NewStatusHandler() *StatusHandler {
	return &StatusHandler{}
}

func (h *StatusHandler) GetStatus(c *fiber.Ctx) error {
	uptime := config.Metrics.GetUptime()
	requestCount := config.Metrics.GetRequestCount()

	return c.JSON(fiber.Map{
		"uptime":        uptime.String(),
		"uptime_seconds": int64(uptime.Seconds()),
		"request_count": requestCount,
		"status":        "ok",
	})
}

