// Package domain ...
package domain

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
)

// User struct
type User struct {
	ID          uint      `gorm:"primaryKey" json:"id,omitempty"`
	ExternalID  string    `json:"external_id" gorm:"uniqueIndex" validate:"required,uuid"`
	Name        string    `json:"name" validate:"required,min=2,max=100"`
	Email       string    `json:"email" validate:"required,email"`
	DateOfBirth time.Time `json:"date_of_birth" validate:"required"`
}

// Custom errors
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidInput      = errors.New("invalid input")
	ErrUserAlreadyExists = errors.New("user already exists")
)

// Validate validates the user struct
func (u *User) Validate() error {
	validate := validator.New()
	if err := validate.Struct(u); err != nil {
		return ErrInvalidInput
	}
	return nil
}
