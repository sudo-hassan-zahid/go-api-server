package database

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/sudo-hassan-zahid/go-api-server/internal/models"
)

func Migrate(db *gorm.DB) error {
	fmt.Println("Running database migrations...")

	err := db.AutoMigrate(
		&models.User{},
	)
	if err != nil {
		return fmt.Errorf("auto migrate: %w", err)
	}

	fmt.Println("Migrations completed successfully")
	return nil
}

func Seed(db *gorm.DB) error {
	fmt.Println("Seeding database...")

	var count int64
	db.Model(&models.User{}).Count(&count)

	if count > 0 {
		fmt.Println("Database already seeded, skipping")
		return nil
	}

	users := []models.User{
		{
			Email:      "admin@example.com",
			Password:   "admin123",
			FirstName:  "Admin",
			LastName:   "User",
			IsVerified: true,
			Role:       "admin",
		},
		{
			Email:      "user@example.com",
			Password:   "user123",
			FirstName:  "Test",
			LastName:   "User",
			IsVerified: true,
			Role:       "user",
		},
	}

	for i := range users {
		if err := db.Create(&users[i]).Error; err != nil {
			return fmt.Errorf("seed user %d: %w", i, err)
		}
	}

	fmt.Printf("Seeded %d users\n", len(users))
	return nil
}

func Setup(db *gorm.DB, shouldSeed bool) error {
	_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := Migrate(db); err != nil {
		return err
	}

	if shouldSeed {
		if err := Seed(db); err != nil {
			return err
		}
	}

	return nil
}
