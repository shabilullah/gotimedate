package router

import (
	"gotimedate/config"
	_ "gotimedate/docs"
	"gotimedate/handlers"
	"gotimedate/middleware"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/gofiber/websocket/v2"
)

func SetupRouter(cfg *config.Config) *fiber.App {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: false,
		Prefork:               cfg.Prefork,
		ErrorHandler:          middleware.ErrorHandler,
	})

	app.Use(recover.New())
	app.Use(middleware.SecurityHeaders())
	app.Use(middleware.Logger(cfg.LogLevel))
	app.Use(middleware.Core(cfg.LogLevel))
	app.Use(etag.New())

	app.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			return cfg.IsOriginAllowed(origin)
		},
		AllowMethods:     strings.Join(cfg.AllowedMethods, ","),
		AllowHeaders:     strings.Join(cfg.AllowedHeaders, ","),
		AllowCredentials: cfg.AllowCredentials,
		MaxAge:           cfg.MaxAge,
	}))

	timeHandler := handlers.NewTimeHandler()
	wsHandler := handlers.NewWSHandler(cfg)

	app.Use("/ws/time", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/time", websocket.New(wsHandler.ServeHTTP, websocket.Config{
		Origins: []string{"*"},
	}))

	app.Get("/swagger/*", swagger.HandlerDefault)

	app.Get("/ws-test", func(c *fiber.Ctx) error {
		indexFile := filepath.Join(cfg.StaticDir, "websocket-test.html")
		return c.SendFile(indexFile)
	})

	app.Get("/health", timeHandler.HealthCheck)

	api := app.Group("/api/v1")
	api.Get("/time", timeHandler.GetCurrentTime)
	api.Get("/timezones", timeHandler.GetAvailableTimezones)
	api.Get("/time/*", timeHandler.GetTimeByTimezone)
	api.Post("/time/convert", timeHandler.ConvertTime)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"name":        "Go TimeDate API",
			"version":     "1.0.0",
			"description": "API for time operations and WebSockets",
		})
	})

	return app
}
