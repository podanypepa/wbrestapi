package repository

import (
	"time"

	"github.com/podanypepa/wbrestapi/internal/domain"
)

// UserEntity represents the GORM model for database
type UserEntity struct {
	ID          uint      `gorm:"primaryKey"`
	ExternalID  string    `gorm:"uniqueIndex"`
	Name        string
	Email       string
	DateOfBirth time.Time
}

// TableName overrides the table name used by UserEntity to `users`
func (UserEntity) TableName() string {
	return "users"
}

// FromDomain converts domain.User entity to UserEntity GORM model
func FromDomain(u *domain.User) *UserEntity {
	return &UserEntity{
		ID:          u.ID,
		ExternalID:  u.ExternalID,
		Name:        u.Name,
		Email:       u.Email,
		DateOfBirth: u.DateOfBirth,
	}
}

// ToDomain converts UserEntity GORM model to domain.User entity
func (e *UserEntity) ToDomain() *domain.User {
	return &domain.User{
		ID:          e.ID,
		ExternalID:  e.ExternalID,
		Name:        e.Name,
		Email:       e.Email,
		DateOfBirth: e.DateOfBirth,
	}
}
