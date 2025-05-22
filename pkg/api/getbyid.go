package api

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/podanypepa/wbrestapi/pkg/repository"
	"gorm.io/gorm"
)

func getByID(userRepository UserRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		externalID := c.Params("id")
		user := repository.User{
			ExternalID: externalID,
		}

		_, err := userRepository.First(&user, externalID)
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
	}
}
