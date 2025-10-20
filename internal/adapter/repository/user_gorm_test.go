package repository

import (
	"testing"
	"time"

	"github.com/podanypepa/wbrestapi/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&domain.User{})
	require.NoError(t, err)

	return db
}

func TestUserGormRepository_Save_Success(t *testing.T) {
	db := setupTestDB(t)
	repo := &UserGormRepository{DB: db}

	user := &domain.User{
		ExternalID:  "550e8400-e29b-41d4-a716-446655440000",
		Name:        "John Doe",
		Email:       "john@example.com",
		DateOfBirth: time.Now(),
	}

	err := repo.Save(user)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)
}

func TestUserGormRepository_Save_DuplicateExternalID(t *testing.T) {
	db := setupTestDB(t)
	repo := &UserGormRepository{DB: db}

	user1 := &domain.User{
		ExternalID:  "550e8400-e29b-41d4-a716-446655440000",
		Name:        "John Doe",
		Email:       "john@example.com",
		DateOfBirth: time.Now(),
	}

	err := repo.Save(user1)
	require.NoError(t, err)

	user2 := &domain.User{
		ExternalID:  "550e8400-e29b-41d4-a716-446655440000",
		Name:        "Jane Doe",
		Email:       "jane@example.com",
		DateOfBirth: time.Now(),
	}

	err = repo.Save(user2)
	assert.Equal(t, domain.ErrUserAlreadyExists, err)
}

func TestUserGormRepository_FindByExternalID_Success(t *testing.T) {
	db := setupTestDB(t)
	repo := &UserGormRepository{DB: db}

	externalID := "550e8400-e29b-41d4-a716-446655440000"
	user := &domain.User{
		ExternalID:  externalID,
		Name:        "John Doe",
		Email:       "john@example.com",
		DateOfBirth: time.Now(),
	}

	err := repo.Save(user)
	require.NoError(t, err)

	found, err := repo.FindByExternalID(externalID)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, externalID, found.ExternalID)
	assert.Equal(t, "John Doe", found.Name)
}

func TestUserGormRepository_FindByExternalID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := &UserGormRepository{DB: db}

	found, err := repo.FindByExternalID("non-existent-uuid")
	assert.Equal(t, domain.ErrUserNotFound, err)
	assert.Nil(t, found)
}
