package usecase

import (
	"github.com/podanypepa/wbrestapi/internal/application/port"
	"github.com/podanypepa/wbrestapi/internal/domain"
)

// SaveUserUseCase ...
type SaveUserUseCase struct {
	Repo port.UserRepository
}

// Execute ...
func (uc *SaveUserUseCase) Execute(user *domain.User) error {
	return uc.Repo.Save(user)
}
