package routes

import (
	"github.com/sudo-hassan-zahid/go-api-server/internal/handler"
	"github.com/sudo-hassan-zahid/go-api-server/internal/repository"
	"github.com/sudo-hassan-zahid/go-api-server/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/gorm"
)

func Setup(app *fiber.App, db *gorm.DB) {
	// CORS configuration
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Initialize layers
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, db)
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
