package routes

import (
	"github.com/sudo-hassan-zahid/go-api-server/internal/handler"
	"github.com/sudo-hassan-zahid/go-api-server/internal/repository"
	"github.com/sudo-hassan-zahid/go-api-server/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Setup(app *fiber.App, db *gorm.DB) {
	// Initialize layers
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// Public routes
	app.Get("/health", handler.HealthCheck)
	api := app.Group("/api")

	// User routes
	api.Post("/users", userHandler.CreateUser)
	api.Post("/users/login", userHandler.LoginUser)
	api.Get("/users", userHandler.GetAllUsers)
	api.Get("/users/:id", userHandler.GetUserByID)

}
