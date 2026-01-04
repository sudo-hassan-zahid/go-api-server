package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/sudo-hassan-zahid/go-api-server/utils"
	"gorm.io/gorm"
)

type User struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Email      string         `gorm:"uniqueIndex;not null" json:"email"`
	Password   string         `gorm:"not null" json:"-"`
	FirstName  string         `gorm:"index:idx_name;not null" json:"first_name"`
	LastName   string         `gorm:"index:idx_name;not null" json:"last_name"`
	IsVerified bool           `gorm:"default:false" json:"is_verified"`
	Role       string         `gorm:"default:'user'" json:"role"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	u.Password, err = utils.HashPassword(u.Password)
	return
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	if tx.Statement.Changed("Password") {
		u.Password, err = utils.HashPassword(u.Password)
	}
	return
}
