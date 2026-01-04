package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	log.Errorf("Error: %s | Path: %s | Method: %s | IP: %s",
		err.Error(),
		c.Path(),
		c.Method(),
		c.IP(),
	)

	return c.Status(code).JSON(fiber.Map{
		"error":   true,
		"message": err.Error(),
		"code":    code,
		"path":    c.Path(),
		"method":  c.Method(),
	})
}

func Logger(logLevel string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if logLevel != "info" && logLevel != "debug" {
			return c.Next()
		}
		start := time.Now()
		err := c.Next()
		log.Infof("%s %s %v", c.Method(), c.Path(), time.Since(start))
		return err
	}
}

func SecurityHeaders() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "DENY")
		c.Set("X-XSS-Protection", "1; mode=block")
		return c.Next()
	}
}

func Core(logLevel string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("requestID", c.Get("X-Request-ID"))
		return c.Next()
	}
}
