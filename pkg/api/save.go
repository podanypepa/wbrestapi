package api

import (
	"cmp"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/podanypepa/wbrestapi/pkg/repository"
)

func save(userRepository UserRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
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

		if _, err := userRepository.Create(&user); err != nil {
			return c.
				Status(fiber.StatusInternalServerError).
				JSON(fiber.Map{
					"error": err.Error(),
				})
		}

		return c.
			Status(fiber.StatusCreated).
			JSON(user)
	}
}
