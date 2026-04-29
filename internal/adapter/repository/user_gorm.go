// Package repository ...
package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mattn/go-sqlite3"
	"github.com/podanypepa/wbrestapi/internal/domain"
	"gorm.io/gorm"
)

// UserGormRepository struct
type UserGormRepository struct {
	DB *gorm.DB
}

// Save ...
func (r *UserGormRepository) Save(ctx context.Context, user *domain.User) error {
	entity := FromDomain(user)
	err := r.DB.WithContext(ctx).Create(entity).Error
	if err != nil {
		// 1. Check for generic GORM duplicated key error
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return domain.ErrUserAlreadyExists
		}

		// 2. Check for PostgreSQL specific unique constraint violation (code 23505)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrUserAlreadyExists
		}

		// 3. Check for SQLite specific unique constraint violation (code 2067)
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && (sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique || sqliteErr.ExtendedCode == 2067) {
			return domain.ErrUserAlreadyExists
		}

		return err
	}
	// Update domain model with generated ID
	user.ID = entity.ID
	return nil
}

// FindByExternalID ...
func (r *UserGormRepository) FindByExternalID(ctx context.Context, externalID string) (*domain.User, error) {
	var entity UserEntity
	if err := r.DB.WithContext(ctx).Where("external_id = ?", externalID).First(&entity).Error; err != nil {
		// Convert GORM specific errors to domain errors
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return entity.ToDomain(), nil
}
