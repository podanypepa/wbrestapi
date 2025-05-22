package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/gofiber/fiber/v2"
	"github.com/podanypepa/wbrestapi/pkg/repository"
)

func TestGetByID(t *testing.T) {
	t.Run("GetByID_OK", func(t *testing.T) {
		userRepository, _ := repository.NewUserRepositoryMock()
		server := NewServer(Config{
			UserRepository: userRepository,
		})

		req, _ := http.NewRequest("GET", fmt.Sprintf("/%s", repository.UserMock.ExternalID), nil)
		res, err := server.Test(req, -1)
		assert.NoError(t, err)

		assert.Equal(t, fiber.StatusOK, res.StatusCode)

		bodyBytes, err := io.ReadAll(res.Body)
		assert.NoError(t, err)

		var apiUser repository.User
		err = json.Unmarshal(bodyBytes, &apiUser)
		assert.NoError(t, err)
		assert.Equal(t, repository.UserMock.ExternalID, apiUser.ExternalID)
	})

	t.Run("GetByID_NOT_FOUND", func(t *testing.T) {
		userRepository, _ := repository.NewUserRepositoryMock()
		server := NewServer(Config{
			UserRepository: userRepository,
		})

		req, _ := http.NewRequest("GET", fmt.Sprintf("/%s", "1"), nil)
		res, err := server.Test(req, -1)
		assert.NoError(t, err)

		assert.Equal(t, fiber.StatusNotFound, res.StatusCode)
	})

}
