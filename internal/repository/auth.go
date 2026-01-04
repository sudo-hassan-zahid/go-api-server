package repository

import (
	"github.com/sudo-hassan-zahid/go-api-server/internal/models"
	"gorm.io/gorm"
)

type AuthRepository interface {
	Create(user *models.User) error
}

type authRepo struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepo{db: db}
}

func (r *authRepo) Create(user *models.User) error {
	return r.db.Create(user).Error
}
