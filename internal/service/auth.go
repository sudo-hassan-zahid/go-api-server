package service

import (
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
	repo     repository.AuthRepository
	userRepo repository.UserRepository
	db       *gorm.DB
}

func NewAuthService(repo repository.AuthRepository, db *gorm.DB) AuthService {
	return &authService{repo: repo, db: db}
}

func (s *authService) CreateUser(email, password, firstName, lastName string) (*models.User, error) {
	var user *models.User
	err := s.db.Transaction(func(tx *gorm.DB) error {
		u := &models.User{
			Email:     email,
			Password:  password,
			FirstName: firstName,
			LastName:  lastName,
		}

		if err := tx.Create(u).Error; err != nil {
			return err
		}

		user = u
		return nil
	})
	return user, err
}

func (s *authService) LoginUser(email, password string) (*models.User, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, appErrors.ErrInvalidCredentials
	}
	if ok := utils.CheckPassword(user.Password, password); !ok {
		return nil, appErrors.ErrInvalidCredentials
	}
	return user, nil
}
