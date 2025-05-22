package api

import (
	"io"
	"net/http"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/gofiber/fiber/v2"
	"github.com/podanypepa/wbrestapi/pkg/repository"
)

func TestIndex(t *testing.T) {
	t.Run("Index", func(t *testing.T) {
		userRepository, _ := repository.NewUserRepositoryMock()
		server := NewServer(Config{
			UserRepository: userRepository,
		})

		req, _ := http.NewRequest("GET", "/", nil)
		res, err := server.Test(req, -1)
		assert.NoError(t, err)

		assert.Equal(t, fiber.StatusOK, res.StatusCode)

		bodyBytes, err := io.ReadAll(res.Body)
		assert.NoError(t, err)
		assert.Equal(t, "wbrestapi", string(bodyBytes))
	})
}
