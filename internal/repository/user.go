package repository

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/sudo-hassan-zahid/go-api-server/internal/database"
	"github.com/sudo-hassan-zahid/go-api-server/internal/domain"
	"github.com/sudo-hassan-zahid/go-api-server/internal/dto"
	"github.com/sudo-hassan-zahid/go-api-server/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByEmail(email string) (*models.User, error)
	GetByID(id uuid.UUID) (*models.User, error)
	GetAll() ([]models.User, error)
	UpdateUser(id uuid.UUID, req dto.UpdateUserRequest) error
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(user *models.User) error {
	if err := r.db.Create(user).Error; err != nil {
		if database.IsUniqueViolation(err) {
			return domain.ErrUserAlreadyExists
		}
		if database.IsNotNullViolation(err) {
			return domain.ErrNotNullViolation
		}
		if database.IsForeignKeyViolation(err) {
			return domain.ErrForeignKeyViolation
		}
		if database.IsCheckViolation(err) {
			return domain.ErrCheckViolation
		}
		return err
	}
	return nil
}

func (r *userRepo) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if database.IsRecordNotFound(err) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) GetByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		if database.IsRecordNotFound(err) {
			return nil, domain.ErrUserNotFound
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrUserNotFound
		}
		return err
	}

	updates := make(map[string]interface{})

	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Password != "" {
		updates["password"] = req.Password
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
			if database.IsUniqueViolation(err) {
				constraint := database.GetConstraintName(err)
				if constraint == "idx_users_email" || strings.Contains(constraint, "email") {
					return domain.ErrUserAlreadyExists
				}
				return domain.ErrUserAlreadyExists
			}
			if database.IsNotNullViolation(err) {
				return domain.ErrNotNullViolation
			}
			if database.IsForeignKeyViolation(err) {
				return domain.ErrForeignKeyViolation
			}
			if database.IsCheckViolation(err) {
				return domain.ErrCheckViolation
			}
			return err
		}
	}

	return nil
}
