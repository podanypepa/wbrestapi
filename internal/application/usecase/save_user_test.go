package usecase

import (
	"testing"
	"time"

	"github.com/podanypepa/wbrestapi/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Save(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByExternalID(externalID string) (*domain.User, error) {
	args := m.Called(externalID)
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
		DateOfBirth: time.Now(),
	}

	mockRepo.On("Save", user).Return(nil)

	err := uc.Execute(user)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestSaveUserUseCase_ValidationError(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := &SaveUserUseCase{Repo: mockRepo}

	user := &domain.User{
		ExternalID:  "invalid-uuid",
		Name:        "John Doe",
		Email:       "john@example.com",
		DateOfBirth: time.Now(),
	}

	err := uc.Execute(user)
	assert.Equal(t, domain.ErrInvalidInput, err)
	mockRepo.AssertNotCalled(t, "Save")
}
