package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/podanypepa/wbrestapi/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Save(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByExternalID(ctx context.Context, externalID string) (*domain.User, error) {
	args := m.Called(ctx, externalID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func TestSaveUserUseCase_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := &SaveUserUseCase{Repo: mockRepo}

	user := &domain.User{
		ExternalID:  "550e8400-e29b-41d4-a716-446655440000",
		Name:        "John Doe",
		Email:       "john@example.com",
		DateOfBirth: time.Now().AddDate(-25, 0, 0),
	}

	ctx := context.Background()
	mockRepo.On("Save", ctx, user).Return(nil)

	err := uc.Execute(ctx, user)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestSaveUserUseCase_InvalidAge(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := &SaveUserUseCase{Repo: mockRepo}

	// User is only 5 years old
	user := &domain.User{
		ExternalID:  "550e8400-e29b-41d4-a716-446655440000",
		Name:        "Young Kid",
		Email:       "kid@example.com",
		DateOfBirth: time.Now().AddDate(-5, 0, 0),
	}

	ctx := context.Background()
	err := uc.Execute(ctx, user)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrInvalidAge, err)
	mockRepo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
}
