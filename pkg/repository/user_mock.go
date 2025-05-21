package repository

import (
	"time"

	"gorm.io/gorm"
)

// UserRepositoryMock ...
type UserRepositoryMock struct{}

// UserMock ...
var UserMock = User{
	ID:          1_000_000,
	ExternalID:  "48904E4A-6B76-4D47-8A96-4A51474179B5",
	Name:        "John Doe",
	Email:       "john@example.com",
	DateOfBirth: time.Now(),
}

// NewUserRepositoryMock ...
func NewUserRepositoryMock() (*UserRepositoryMock, error) {
	return &UserRepositoryMock{}, nil
}

// Create (save) User to DB
// returns ID od new User and error
func (u *UserRepositoryMock) Create(newUser *User) (int, error) {
	return int(newUser.ID), nil
}

// First return first User by UUID
func (u *UserRepositoryMock) First(_ *User, externalID string) (*User, error) {
	if externalID == UserMock.ExternalID {
		return &UserMock, nil
	}
	return nil, gorm.ErrRecordNotFound
}

// DB ...
func (u *UserRepositoryMock) DB() *gorm.DB {
	return nil
}
