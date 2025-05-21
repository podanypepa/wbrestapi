// Package repository ...
package repository

import (
	"time"

	"gorm.io/gorm"
)

// UserRepository ...
type UserRepository struct {
	db     *gorm.DB
	config UserRepositoryConfig
}

// User model
type User struct {
	ID          uint      `gorm:"primaryKey"`
	ExternalID  string    `json:"external_id" gorm:"uniqueIndex"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	DateOfBirth time.Time `json:"date_of_birth"`
}

// UserRepositoryConfig ...
type UserRepositoryConfig struct {
	Db *gorm.DB
}

// NewUserRepository create new user repository
func NewUserRepository(cfg UserRepositoryConfig) (*UserRepository, error) {
	u := &UserRepository{
		db:     cfg.Db,
		config: cfg,
	}
	return u, nil
}

// Create (save) User to DB
// returns ID od new User and error
func (u *UserRepository) Create(newUser *User) (int, error) {
	res := u.db.Create(newUser)
	return int(newUser.ID), res.Error
}

// First return first User by UUID
func (u *UserRepository) First(dbUser *User, externalID string) (*User, error) {
	dbUser.ExternalID = externalID
	res := u.db.First(dbUser, "external_id = ?", externalID)
	return dbUser, res.Error
}

// DB ...
func (u *UserRepository) DB() *gorm.DB {
	return u.db
}
