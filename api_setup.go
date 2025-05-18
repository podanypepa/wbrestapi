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

	return app
}
