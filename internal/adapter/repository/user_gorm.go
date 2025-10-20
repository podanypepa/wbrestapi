// Package repository ...
package repository

import (
	"errors"

	"github.com/podanypepa/wbrestapi/internal/domain"
	"gorm.io/gorm"
)

// UserGormRepository struct
type UserGormRepository struct {
	DB *gorm.DB
}

// Save ...
func (r *UserGormRepository) Save(user *domain.User) error {
	err := r.DB.Create(user).Error
	if err != nil {
		// Convert GORM specific errors to domain errors
		errMsg := err.Error()
		// Check for unique constraint violations (PostgreSQL and SQLite)
		if errors.Is(err, gorm.ErrDuplicatedKey) ||
			(errMsg != "" && (
				// PostgreSQL: "duplicate key value violates unique constraint"
				// SQLite: "UNIQUE constraint failed"
				errors.Is(err, gorm.ErrDuplicatedKey) ||
				len(errMsg) > 15 && errMsg[:15] == "UNIQUE constrai")) {
			return domain.ErrUserAlreadyExists
		}
		return err
	}
	return nil
}

// FindByExternalID ...
func (r *UserGormRepository) FindByExternalID(externalID string) (*domain.User, error) {
	var user domain.User
	if err := r.DB.Where("external_id = ?", externalID).First(&user).Error; err != nil {
		// Convert GORM specific errors to domain errors
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}
