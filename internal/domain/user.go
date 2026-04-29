// Package domain ...
package domain

import (
	"errors"
	"time"
)

// User struct represents the core domain model
type User struct {
	ID          uint
	ExternalID  string
	Name        string
	Email       string
	DateOfBirth time.Time
}

// Custom errors
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidInput      = errors.New("invalid input")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidAge        = errors.New("user must be at least 15 years old")
	ErrFutureBirthDate   = errors.New("date of birth cannot be in the future")
)

// Validate checks business invariants
func (u *User) Validate() error {
	now := time.Now()

	if u.DateOfBirth.After(now) {
		return ErrFutureBirthDate
	}

	// Calculate age
	age := now.Year() - u.DateOfBirth.Year()
	if now.YearDay() < u.DateOfBirth.YearDay() {
		age--
	}

	if age < 15 {
		return ErrInvalidAge
	}

	return nil
}
