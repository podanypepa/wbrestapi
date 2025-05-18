package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetByID(t *testing.T) {
	t.Run("GetByID_OK", func(t *testing.T) {
		mockDb, mock, cleanup := setupMockDB(t)
		defer cleanup()

		db = mockDb
		app := apiSetup()
		dbUser := User{
			ID:          1,
			ExternalID:  "123e4567-e89b-12d3-a456-426614174000",
			Name:        "John Doe",
			Email:       "john@example.com",
			DateOfBirth: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		rows := sqlmock.NewRows([]string{"id", "external_id", "name", "email", "date_of_birth"}).
			AddRow(dbUser.ID, dbUser.ExternalID, dbUser.Name, dbUser.Email, dbUser.DateOfBirth)

		mock.ExpectQuery(`SELECT \* FROM \"users\" WHERE external_id = \$1 ORDER BY \"users\".\"id\" LIMIT \$2`).
			WithArgs(dbUser.ExternalID, 1).
			WillReturnRows(rows)

		req, _ := http.NewRequest("GET", fmt.Sprintf("/%s", dbUser.ExternalID), nil)

		res, err := app.Test(req, -1)

		assert.NoError(t, err)
		bodyBytes, err := io.ReadAll(res.Body)
		assert.NoError(t, err)

		var apiUser User
		err = json.Unmarshal(bodyBytes, &apiUser)
		assert.NoError(t, err)
		assert.Equal(t, dbUser.ExternalID, apiUser.ExternalID)
	})
}
