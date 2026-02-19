package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/sudo-hassan-zahid/go-api-server/internal/dto"
	customErr "github.com/sudo-hassan-zahid/go-api-server/internal/errors"
	customerrors "github.com/sudo-hassan-zahid/go-api-server/internal/errors"
	"github.com/sudo-hassan-zahid/go-api-server/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByEmail(email string) (*models.User, error)
	GetByID(id uuid.UUID) (*models.User, error)
	GetAll() ([]models.User, error)
	UpdateUser(id uuid.UUID, req dto.UpdateUserRequest) error
	DeleteUser(id uuid.UUID) error
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(user *models.User) error {
	if err := r.db.Create(user).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return customErr.ErrUserAlreadyExists
		}
		return err
	}
	return nil
}

func (r *userRepo) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, customErr.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) GetByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, customErr.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) GetAll() ([]models.User, error) {
	var users []models.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepo) UpdateUser(id uuid.UUID, req dto.UpdateUserRequest) error {
	var user models.User

	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return customerrors.ParseDatabaseError(err)
	}

	updates := make(map[string]interface{})

	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.FirstName != "" {
		updates["first_name"] = req.FirstName
	}
	if req.LastName != "" {
		updates["last_name"] = req.LastName
	}
	if req.Role != "" {
		updates["role"] = req.Role
	}

	if len(updates) > 0 {
		if err := r.db.Model(&user).Updates(updates).Error; err != nil {
			return customerrors.ParseDatabaseError(err)
		}
	}

	return nil
}

func (r *userRepo) DeleteUser(id uuid.UUID) error {
	result := r.db.Delete(&models.User{}, "id = ?", id)

	if result.Error != nil {
		return customerrors.ParseDatabaseError(result.Error)
	}

	if result.RowsAffected == 0 {
		return customerrors.ErrUserNotFound
	}

	return nil
}
