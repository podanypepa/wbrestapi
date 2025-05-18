package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	assert.NoError(t, err)

	cleanup := func() {
		db.Close()
	}

	return gormDB, mock, cleanup
}

func TestCreateUser(t *testing.T) {
	mockDb, mock, cleanup := setupMockDB(t)
	defer cleanup()

	db = mockDb
	app := apiSetup()

	t.Run("CreateUser", func(t *testing.T) {
		user := &User{
			ExternalID:  uuid.NewString(),
			Name:        "Jane Doe",
			Email:       "jane@example.com",
			DateOfBirth: time.Date(2025, 5, 18, 13, 10, 50, 801129000, time.UTC),
		}

		mock.ExpectBegin()
		mock.
			ExpectQuery(
				regexp.QuoteMeta(`INSERT INTO "users" ("external_id","name","email","date_of_birth") VALUES ($1,$2,$3,$4) RETURNING "id"`),
			).
			WithArgs(user.ExternalID, user.Name, user.Email, user.DateOfBirth).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		b, _ := json.Marshal(user)
		req, _ := http.NewRequest("POST", "/save", strings.NewReader(string(b)))
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Content-Length", fmt.Sprintf("%d", utf8.RuneCountInString(string(b))))

		res, err := app.Test(req, -1)

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

		res, err := app.Test(req, -1)
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
