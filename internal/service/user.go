package service

import (
	appErrors "github.com/sudo-hassan-zahid/go-api-server/internal/errors"

	"github.com/sudo-hassan-zahid/go-api-server/internal/models"
	"github.com/sudo-hassan-zahid/go-api-server/internal/repository"
	"gorm.io/gorm"
)

type UserService interface {
	GetAllUsers() ([]models.User, error)
	GetUserByID(id uint) (*models.User, error)
}

type userService struct {
	repo repository.UserRepository
	db   *gorm.DB
}

func NewUserService(repo repository.UserRepository, db *gorm.DB) UserService {
	return &userService{repo: repo, db: db}
}

func (s *userService) GetAllUsers() ([]models.User, error) {
	return s.repo.GetAll()
}

func (s *userService) GetUserByID(id uint) (*models.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, appErrors.ErrUserNotFound
	}
	return user, nil
}
