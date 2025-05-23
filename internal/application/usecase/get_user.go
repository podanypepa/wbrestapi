// Package usecase ...
package usecase

import (
	"github.com/podanypepa/wbrestapi/internal/application/port"
	"github.com/podanypepa/wbrestapi/internal/domain"
)

// GetUserUseCase struct
type GetUserUseCase struct {
	Repo port.UserRepository
}

// Execute ...
func (uc *GetUserUseCase) Execute(externalID string) (*domain.User, error) {
	return uc.Repo.FindByExternalID(externalID)
}
