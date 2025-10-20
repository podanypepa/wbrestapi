package usecase

import (
	"github.com/podanypepa/wbrestapi/internal/application/port"
	"github.com/podanypepa/wbrestapi/internal/domain"
)

// SaveUserUseCase ...
type SaveUserUseCase struct {
	Repo port.UserRepository
}

// Execute validates and saves user
func (uc *SaveUserUseCase) Execute(user *domain.User) error {
	if err := user.Validate(); err != nil {
		return err
	}
	return uc.Repo.Save(user)
}
