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
)
