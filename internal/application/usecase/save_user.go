package usecase

import (
	"github.com/podanypepa/wbrestapi/internal/application/port"
	"github.com/podanypepa/wbrestapi/internal/domain"
)

type SaveUserUseCase struct {
	Repo port.UserRepository
}

func (uc *SaveUserUseCase) Execute(user *domain.User) error {
	return uc.Repo.Save(user)
}
