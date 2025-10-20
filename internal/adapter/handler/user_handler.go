// Package handler ...
package handler

import (
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/podanypepa/wbrestapi/internal/application/port"
	"github.com/podanypepa/wbrestapi/internal/domain"
)

// UserHandler struct
type UserHandler struct {
	SaveUC port.SaveUserExecutor
	GetUC  port.GetUserExecutor
	Logger *slog.Logger
}

// RegisterRoutes ...
func (h *UserHandler) RegisterRoutes(app *fiber.App) {
	app.Get("/healthz", h.HealthCheck)
	app.Post("/save", h.SaveUser)
	app.Get("/:id", h.GetUser)
}

// SaveUser ...
func (h *UserHandler) SaveUser(c *fiber.Ctx) error {
	var user domain.User
	if err := c.BodyParser(&user); err != nil {
		h.Logger.Error("failed to parse request body", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid payload",
		})
	}

	if err := h.SaveUC.Execute(&user); err != nil {
		if errors.Is(err, domain.ErrInvalidInput) {
			h.Logger.Warn("invalid user input", "error", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "validation failed: " + err.Error(),
			})
		}
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			h.Logger.Warn("user already exists", "external_id", user.ExternalID)
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "user with this external_id already exists",
			})
		}
		h.Logger.Error("failed to save user", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	h.Logger.Info("user saved successfully", "external_id", user.ExternalID)
	return c.Status(fiber.StatusCreated).JSON(user)
}

// GetUser ...
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	externalID := c.Params("id")

	user, err := h.GetUC.Execute(externalID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			h.Logger.Info("user not found", "external_id", externalID)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "user not found",
			})
		}
		h.Logger.Error("failed to get user", "error", err, "external_id", externalID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	return c.JSON(user)
}

// HealthCheck ...
func (h *UserHandler) HealthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
}
