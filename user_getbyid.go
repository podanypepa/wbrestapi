package main

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// GET /:id
func getUserByID(c *fiber.Ctx) error {
	externalID := c.Params("id")

	user := User{
		ExternalID: externalID,
	}

	res := db.First(&user, "external_id = ?", externalID)
	if res.Error != nil {
		switch {
		case errors.Is(res.Error, gorm.ErrRecordNotFound):
			return c.
				Status(fiber.StatusNotFound).
				JSON(apiError{Error: fmt.Sprintf("user %s not found", externalID)})
		default:
			return c.
				Status(fiber.StatusInternalServerError).
				JSON(apiError{Error: res.Error.Error()})
		}
	}

	return c.JSON(user)
}
