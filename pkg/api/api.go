// Package api ...
package api

import (
	"cmp"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
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

	app.Get("/healthz", healthCheck(cfg.UserRepository.DB()))
	app.Get("/", index())
	app.Post("/save", save(cfg.UserRepository))
	app.Get("/:id", getByID(cfg.UserRepository))

	return app
}
