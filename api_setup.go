package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func apiSetup() *fiber.App {
	app := fiber.New(fiber.Config{
		IdleTimeout:           idleTimeout,
		DisableStartupMessage: true,
	})

	app.Use(logger.New())

	app.Post("/save", createUser)
	app.Get("/:id", getUserByID)

	app.Get("/healthz", func(c *fiber.Ctx) error {
		sqlDB, err := db.DB()
		if err != nil || sqlDB.Ping() != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status": "unhealthy",
				"error":  "database not reachable",
			})
		}

		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	return app
}
