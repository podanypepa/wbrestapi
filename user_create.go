package main

import (
	"cmp"
	"errors"
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// POST /save
func createUser(c *fiber.Ctx) error {
	var user User
	if err := c.BodyParser(&user); err != nil {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(apiError{
				Error: fmt.Sprintf("cannot parse JSON: %s", err.Error()),
			})
	}

	user.ExternalID = cmp.Or(user.ExternalID, uuid.New().String())
	if _, err := uuid.Parse(user.ExternalID); err != nil {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(apiError{
				Error: fmt.Sprintf("invalid uuid: %s", user.ExternalID),
			})
	}

	res := db.Create(&user)
	if res.Error != nil {
		slog.Error("createUser", "err", res.Error)
		switch {
		case errors.Is(res.Error, gorm.ErrDuplicatedKey):
			return c.
				Status(fiber.StatusNotAcceptable).
				JSON(apiError{
					Error: "Duplicid external_id.",
				})
		default:
			return c.
				Status(fiber.StatusInternalServerError).
				JSON(apiError{
					Error: res.Error.Error(),
				})
		}
	}

	return c.
		Status(fiber.StatusCreated).
		JSON(user)
}
