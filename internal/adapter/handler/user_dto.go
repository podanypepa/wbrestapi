package handler

import (
	"time"

	"github.com/podanypepa/wbrestapi/internal/domain"
)

// UserRequest represents the JSON payload for creating/saving a user
type UserRequest struct {
	ExternalID  string    `json:"external_id" validate:"required,uuid"`
	Name        string    `json:"name" validate:"required,min=2,max=100"`
	Email       string    `json:"email" validate:"required,email"`
	DateOfBirth time.Time `json:"date_of_birth" validate:"required"`
}

// UserResponse represents the JSON response for a user
type UserResponse struct {
	ID          uint      `json:"id"`
	ExternalID  string    `json:"external_id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	DateOfBirth time.Time `json:"date_of_birth"`
}

// ToDomain converts UserRequest DTO to domain.User entity
func (r *UserRequest) ToDomain() *domain.User {
	return &domain.User{
		ExternalID:  r.ExternalID,
		Name:        r.Name,
		Email:       r.Email,
		DateOfBirth: r.DateOfBirth,
	}
}

// NewUserResponse creates a UserResponse DTO from a domain.User entity
func NewUserResponse(u *domain.User) *UserResponse {
	return &UserResponse{
		ID:          u.ID,
		ExternalID:  u.ExternalID,
		Name:        u.Name,
		Email:       u.Email,
		DateOfBirth: u.DateOfBirth,
	}
}
