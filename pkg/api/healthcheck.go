package api

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// HealthCheckResponse ...
type HealthCheckResponse struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

func healthCheck(gormDB *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sqlDB, err := gormDB.DB()

		if err != nil || sqlDB.Ping() != nil {
			return c.
				Status(fiber.StatusServiceUnavailable).
				JSON(HealthCheckResponse{
					Status: "unhealthy",
					Error:  "database not reachable",
				})
		}

		return c.JSON(HealthCheckResponse{
			Status: "ok",
		})
	}
}
