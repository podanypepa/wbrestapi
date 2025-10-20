package usecase

import (
	"testing"
	"time"

	"github.com/podanypepa/wbrestapi/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestGetUserUseCase_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := &GetUserUseCase{Repo: mockRepo}

	externalID := "550e8400-e29b-41d4-a716-446655440000"
	expectedUser := &domain.User{
		ID:          1,
		ExternalID:  externalID,
		Name:        "John Doe",
		Email:       "john@example.com",
		DateOfBirth: time.Now(),
	}

	mockRepo.On("FindByExternalID", externalID).Return(expectedUser, nil)

	user, err := uc.Execute(externalID)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockRepo.AssertExpectations(t)
}

func TestGetUserUseCase_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := &GetUserUseCase{Repo: mockRepo}

	externalID := "non-existent"
	mockRepo.On("FindByExternalID", externalID).Return(nil, domain.ErrUserNotFound)

	user, err := uc.Execute(externalID)
	assert.Equal(t, domain.ErrUserNotFound, err)
	assert.Nil(t, user)
	mockRepo.AssertExpectations(t)
}
