package usecase

import (
	"context"

	"github.com/podanypepa/wbrestapi/internal/application/port"
	"github.com/podanypepa/wbrestapi/internal/domain"
)

// SaveUserUseCase is the use case for saving a user
type SaveUserUseCase struct {
	Repo port.UserRepository
}

// Execute saves user
func (uc *SaveUserUseCase) Execute(ctx context.Context, user *domain.User) error {
	if err := user.Validate(); err != nil {
		return err
	}
	return uc.Repo.Save(ctx, user)
}
