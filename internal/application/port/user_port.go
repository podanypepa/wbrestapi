// Package port ...
package port

import (
	"context"

	"github.com/podanypepa/wbrestapi/internal/domain"
)

// SaveUserExecutor interface
type SaveUserExecutor interface {
	Execute(ctx context.Context, user *domain.User) error
}

// GetUserExecutor interface
type GetUserExecutor interface {
	Execute(ctx context.Context, externalID string) (*domain.User, error)
}

// UserRepository interface
type UserRepository interface {
	Save(ctx context.Context, user *domain.User) error
	FindByExternalID(ctx context.Context, externalID string) (*domain.User, error)
}
