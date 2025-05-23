package repository

import (
	"github.com/podanypepa/wbrestapi/internal/domain"
	"gorm.io/gorm"
)

type UserGormRepository struct {
	DB *gorm.DB
}

func (r *UserGormRepository) Save(user *domain.User) error {
	return r.DB.Create(user).Error
}

func (r *UserGormRepository) FindByExternalID(externalID string) (*domain.User, error) {
	var user domain.User
	if err := r.DB.Where("external_id = ?", externalID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
