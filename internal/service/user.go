package service

import (
	"github.com/google/uuid"
	"github.com/sudo-hassan-zahid/go-api-server/internal/dto"
	customErr "github.com/sudo-hassan-zahid/go-api-server/internal/errors"
	"github.com/sudo-hassan-zahid/go-api-server/internal/models"
	"github.com/sudo-hassan-zahid/go-api-server/internal/repository"
	"gorm.io/gorm"
)

type UserService interface {
	GetAllUsers() ([]models.User, error)
	GetUserByID(id uuid.UUID) (*models.User, error)
	UpdateUser(id uuid.UUID, req dto.UpdateUserRequest) error
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

func (s *userService) GetUserByID(id uuid.UUID) (*models.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, customErr.ErrUserNotFound
	}
	return user, nil
}

func (s *userService) UpdateUser(id uuid.UUID, req dto.UpdateUserRequest) error {
	err := s.repo.UpdateUser(id, req)
	if err != nil {
		return err
	}
	return nil
}
