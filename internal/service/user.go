package service

import (
	"errors"

	"github.com/sudo-hassan-zahid/go-api-server/internal/models"
	"github.com/sudo-hassan-zahid/go-api-server/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	CreateUser(email, password, firstName, lastName string) (*models.User, error)
	LoginUser(email, password string) (*models.User, error)
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

func (s *userService) CreateUser(email, password, firstName, lastName string) (*models.User, error) {
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

func (s *userService) LoginUser(email, password string) (*models.User, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}
	return user, nil
}

func (s *userService) GetAllUsers() ([]models.User, error) {
	return s.repo.GetAll()
}

func (s *userService) GetUserByID(id uint) (*models.User, error) {
	return s.repo.GetByID(id)
}
