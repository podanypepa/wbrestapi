package usecase

import (
	"context"
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

	ctx := context.Background()
	mockRepo.On("FindByExternalID", ctx, externalID).Return(expectedUser, nil)

	user, err := uc.Execute(ctx, externalID)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockRepo.AssertExpectations(t)
}

func TestGetUserUseCase_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := &GetUserUseCase{Repo: mockRepo}

	externalID := "non-existent"
	ctx := context.Background()
	mockRepo.On("FindByExternalID", ctx, externalID).Return(nil, domain.ErrUserNotFound)

	user, err := uc.Execute(ctx, externalID)
	assert.Equal(t, domain.ErrUserNotFound, err)
	assert.Nil(t, user)
	mockRepo.AssertExpectations(t)
}
