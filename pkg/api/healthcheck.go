package api

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func healthCheck(gormDB *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sqlDB, err := gormDB.DB()
		if err != nil || sqlDB.Ping() != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status": "unhealthy",
				"error":  "database not reachable",
			})
		}

		return c.JSON(fiber.Map{
			"status": "ok",
		})

	}
}
