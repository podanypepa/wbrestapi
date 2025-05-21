package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/alecthomas/assert"
	"github.com/gofiber/fiber"
	"github.com/google/uuid"
	"github.com/podanypepa/wbrestapi/pkg/api"
	"github.com/podanypepa/wbrestapi/pkg/repository"
)

func TestCreateUser(t *testing.T) {
	t.Run("CreateUser", func(t *testing.T) {
		user := &User{
			ExternalID:  uuid.NewString(),
			Name:        "Jane Doe",
			Email:       "jane@example.com",
			DateOfBirth: time.Date(2025, 5, 18, 13, 10, 50, 801129000, time.UTC),
		}

		b, _ := json.Marshal(user)

		req, _ := http.NewRequest("POST", "/save", strings.NewReader(string(b)))
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Content-Length", fmt.Sprintf("%d", utf8.RuneCountInString(string(b))))

		userRepository, _ := repository.NewUserRepositoryMock()
		server := api.NewServer(api.Config{
			UserRepository: userRepository,
		})

		res, err := server.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, res.StatusCode, fiber.StatusCreated)

		bodyBytes, err := io.ReadAll(res.Body)
		assert.NoError(t, err)

		var createdUser User
		err = json.Unmarshal(bodyBytes, &createdUser)
		assert.NoError(t, err)
		assert.Equal(t, user.ExternalID, createdUser.ExternalID)
	})

	t.Run("InvalidUUID", func(t *testing.T) {
		user := &User{
			ExternalID: "1",
		}

		b, _ := json.Marshal(user)
		req, _ := http.NewRequest("POST", "/save", strings.NewReader(string(b)))
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Content-Length", fmt.Sprintf("%d", utf8.RuneCountInString(string(b))))

		userRepository, _ := repository.NewUserRepositoryMock()
		server := api.NewServer(api.Config{
			UserRepository: userRepository,
		})

		res, err := server.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, res.StatusCode, fiber.StatusBadRequest)

		bodyBytes, err := io.ReadAll(res.Body)
		assert.NoError(t, err)

		var aErr apiError
		err = json.Unmarshal(bodyBytes, &aErr)
		assert.NoError(t, err)

		assert.Equal(t, aErr.Error, fmt.Sprintf("invalid uuid: %s", user.ExternalID))
	})
}
