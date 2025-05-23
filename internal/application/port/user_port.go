package port

import "github.com/podanypepa/wbrestapi/internal/domain"

type SaveUserExecutor interface {
	Execute(user *domain.User) error
}

type GetUserExecutor interface {
	Execute(externalID string) (*domain.User, error)
}

type UserRepository interface {
	Save(user *domain.User) error
	FindByExternalID(externalID string) (*domain.User, error)
}
