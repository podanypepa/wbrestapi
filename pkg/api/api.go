// Package api ...
package api

import (
	"cmp"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/google/uuid"
	"github.com/podanypepa/wbrestapi/pkg/repository"
	"gorm.io/gorm"
)

// Config of API server
type Config struct {
	Port           string
	UserRepository UserRepository
	IdleTimeout    time.Duration
}

const (
	defaultIdleTimeout = 5 * time.Second
)

// UserRepository interface
type UserRepository interface {
	Create(*repository.User) (int, error)
	First(*repository.User, string) (*repository.User, error)
	DB() *gorm.DB
}

// NewServer returns new REST API server
func NewServer(cfg Config) *fiber.App {
	app := fiber.New(fiber.Config{
		IdleTimeout:           cmp.Or(cfg.IdleTimeout, defaultIdleTimeout),
		DisableStartupMessage: true,
	})

	app.Use(logger.New())

	app.Get("/healthz", func(c *fiber.Ctx) error {
		sqlDB, err := cfg.UserRepository.DB().DB()
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

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("wbrestapi")
	})

	app.Post("/save", func(c *fiber.Ctx) error {
		var user repository.User
		if err := c.BodyParser(&user); err != nil {
			return c.
				Status(fiber.StatusBadRequest).
				JSON(fiber.Map{
					"error": fmt.Sprintf("cannot parse JSON: %s", err.Error()),
				})
		}

		user.ExternalID = cmp.Or(user.ExternalID, uuid.New().String())
		if _, err := uuid.Parse(user.ExternalID); err != nil {
			return c.
				Status(fiber.StatusBadRequest).
				JSON(fiber.Map{
					"error": fmt.Sprintf("invalid uuid: %s", user.ExternalID),
				})
		}

		if _, err := cfg.UserRepository.Create(&user); err != nil {
			return c.
				Status(fiber.StatusInternalServerError).
				JSON(fiber.Map{
					"error": err.Error(),
				})
		}

		return c.
			Status(fiber.StatusCreated).
			JSON(user)
	})

	app.Get("/:id", func(c *fiber.Ctx) error {
		externalID := c.Params("id")
		user := repository.User{
			ExternalID: externalID,
		}

		_, err := cfg.UserRepository.First(&user, externalID)
		if err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):
				return c.
					Status(fiber.StatusNotFound).
					JSON(fiber.Map{
						"rrror": fmt.Sprintf("user %s not found", externalID),
					})
			default:
				return c.
					Status(fiber.StatusInternalServerError).
					JSON(fiber.Map{
						"error": err.Error(),
					})
			}
		}
		return c.JSON(user)
	})

	return app
}
