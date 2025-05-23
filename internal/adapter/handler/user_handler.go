// Package handler ...
package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/podanypepa/wbrestapi/internal/application/port"
	"github.com/podanypepa/wbrestapi/internal/domain"
)

// UserHandler struct
type UserHandler struct {
	SaveUC port.SaveUserExecutor
	GetUC  port.GetUserExecutor
}

// RegisterRoutes ...
func (h *UserHandler) RegisterRoutes(app *fiber.App) {
	app.Post("/save", h.SaveUser)
	app.Get("/:id", h.GetUser)
}

// SaveUser ...
func (h *UserHandler) SaveUser(c *fiber.Ctx) error {
	var user domain.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}

	if err := h.SaveUC.Execute(&user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

// GetUser ...
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	externalID := c.Params("id")

	user, err := h.GetUC.Execute(externalID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}
	return c.JSON(user)
}
