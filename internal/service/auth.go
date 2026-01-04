package service

import (
	"github.com/google/uuid"
	appErrors "github.com/sudo-hassan-zahid/go-api-server/internal/errors"
	"github.com/sudo-hassan-zahid/go-api-server/internal/models"
	"github.com/sudo-hassan-zahid/go-api-server/internal/repository"
	"github.com/sudo-hassan-zahid/go-api-server/utils"
	"gorm.io/gorm"
)

type AuthService interface {
	CreateUser(email, password, firstName, lastName string) (*models.User, error)
	LoginUser(email, password string) (*models.User, error)
}

type authService struct {
	repo repository.AuthRepository
	db   *gorm.DB
}

func NewAuthService(repo repository.AuthRepository, db *gorm.DB) AuthService {
	return &authService{repo: repo, db: db}
}

func (s *authService) CreateUser(email, password, firstName, lastName string) (*models.User, error) {
	var exists bool
	if err := s.db.Model(&models.User{}).Select("count(*) > 0").Where("email = ?", email).Find(&exists).Error; err != nil {
		return nil, err
	}
	if exists {
		return nil, appErrors.ErrEmailAlreadyExists
	}

	hashed, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:        uuid.New(),
		Email:     email,
		Password:  hashed,
		FirstName: firstName,
		LastName:  lastName,
		Role:      "user",
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) LoginUser(email, password string) (*models.User, error) {
	user := &models.User{}
	if err := s.db.Where("email = ?", email).First(user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, appErrors.ErrInvalidCredentials
		}
		return nil, err
	}

	if !utils.CheckPassword(user.Password, password) {
		return nil, appErrors.ErrInvalidCredentials
	}

	return user, nil
}
