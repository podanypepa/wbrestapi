package api

import (
	"github.com/gofiber/fiber/v2"
)

func index() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendString("wbrestapi")
	}
}
