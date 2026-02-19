package routes

import (
	"github.com/sudo-hassan-zahid/go-api-server/internal/auth"
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

	// Rate limiter
	authRateLimiter := middleware.AuthRateLimiter()
	publicRateLimiter := middleware.PublicRateLimiter()

	// API group
	api := app.Group("/api")

	// Auth APIs
	authRepo := repository.NewUserRepository(db)
	authService := auth.NewService(authRepo)
	authHandler := handler.NewAuthHandler(authService)
	authRoutes := api.Group("/auth")
	authRoutes.Post("/signup", publicRateLimiter, authHandler.CreateUser)
	authRoutes.Post("/login", publicRateLimiter, authHandler.LoginUser)
	authRoutes.Post("/refresh", publicRateLimiter, authHandler.RefreshToken)
	authRoutes.Post("/logout", jwt, authHandler.Logout)

	// User APIs
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, db)
	userHandler := handler.NewUserHandler(userService)
	users := api.Group("/users")
	users.Get("/", jwt, authRateLimiter, userHandler.GetAllUsers)
	users.Get("/:id", jwt, authRateLimiter, userHandler.GetUserByID)
	users.Patch("/:id", jwt, authRateLimiter, userHandler.UpdateUser)
	users.Delete("/:id", jwt, authRateLimiter, userHandler.DeleteUser)

	// Public routes
	publicHandler := handler.NewPublicHandler(db)
	api.Get("/health/server", publicRateLimiter, publicHandler.HealthCheckServer)
	api.Get("/health/db", publicRateLimiter, publicHandler.HealthCheckDB)
}
