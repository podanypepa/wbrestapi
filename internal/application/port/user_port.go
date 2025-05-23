// Package port ...
package port

import "github.com/podanypepa/wbrestapi/internal/domain"

// SaveUserExecutor interface
type SaveUserExecutor interface {
	Execute(user *domain.User) error
}

// GetUserExecutor interface
type GetUserExecutor interface {
	Execute(externalID string) (*domain.User, error)
}

// UserRepository interface
type UserRepository interface {
	Save(user *domain.User) error
	FindByExternalID(externalID string) (*domain.User, error)
}
