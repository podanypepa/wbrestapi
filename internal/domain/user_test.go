package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUser_Validate_Success(t *testing.T) {
	user := User{
		ExternalID:  "550e8400-e29b-41d4-a716-446655440000",
		Name:        "John Doe",
		Email:       "john@example.com",
		DateOfBirth: time.Now().AddDate(-25, 0, 0),
	}

	err := user.Validate()
	assert.NoError(t, err)
}

func TestUser_Validate_InvalidUUID(t *testing.T) {
	user := User{
		ExternalID:  "not-a-uuid",
		Name:        "John Doe",
		Email:       "john@example.com",
		DateOfBirth: time.Now(),
	}

	err := user.Validate()
	assert.Equal(t, ErrInvalidInput, err)
}

func TestUser_Validate_InvalidEmail(t *testing.T) {
	user := User{
		ExternalID:  "550e8400-e29b-41d4-a716-446655440000",
		Name:        "John Doe",
		Email:       "not-an-email",
		DateOfBirth: time.Now(),
	}

	err := user.Validate()
	assert.Equal(t, ErrInvalidInput, err)
}

func TestUser_Validate_MissingName(t *testing.T) {
	user := User{
		ExternalID:  "550e8400-e29b-41d4-a716-446655440000",
		Name:        "",
		Email:       "john@example.com",
		DateOfBirth: time.Now(),
	}

	err := user.Validate()
	assert.Equal(t, ErrInvalidInput, err)
}

func TestUser_Validate_NameTooShort(t *testing.T) {
	user := User{
		ExternalID:  "550e8400-e29b-41d4-a716-446655440000",
		Name:        "J",
		Email:       "john@example.com",
		DateOfBirth: time.Now(),
	}

	err := user.Validate()
	assert.Equal(t, ErrInvalidInput, err)
}
