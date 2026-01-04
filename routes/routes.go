package routes

import (
	"github.com/sudo-hassan-zahid/go-api-server/internal/handler"
	"github.com/sudo-hassan-zahid/go-api-server/internal/middleware"
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

	// JWT auth
	jwt := middleware.JWTMiddleware()

	// API group
	api := app.Group("/api")

	// Auth APIs
	authRepo := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepo, db)
	authHandler := handler.NewAuthHandler(authService)
	auth := api.Group("/auth")
	auth.Post("/signup", authHandler.CreateUser)
	auth.Post("/login", authHandler.LoginUser)

	// User APIs
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, db)
	userHandler := handler.NewUserHandler(userService)
	users := api.Group("/users")
	users.Get("/", jwt, userHandler.GetAllUsers)
	users.Get("/:id", jwt, userHandler.GetUserByID)

	// Public routes
	publicHandler := handler.NewPublicHandler()
	api.Get("/health", publicHandler.HealthCheck)
}
