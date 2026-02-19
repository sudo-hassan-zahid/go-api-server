package routes

import (
	"github.com/sudo-hassan-zahid/go-api-server/internal/auth"
	"github.com/sudo-hassan-zahid/go-api-server/internal/config"
	"github.com/sudo-hassan-zahid/go-api-server/internal/constants"
	"github.com/sudo-hassan-zahid/go-api-server/internal/handler"
	"github.com/sudo-hassan-zahid/go-api-server/internal/middleware"
	"github.com/sudo-hassan-zahid/go-api-server/internal/repository"
	"github.com/sudo-hassan-zahid/go-api-server/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/gorm"
)

func Setup(app *fiber.App, db *gorm.DB, cfg *config.Config) {
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
	loginRateLimiter := middleware.LoginRateLimiter()

	// API group
	api := app.Group("/api")

	// Auth APIs
	authRepo := repository.NewUserRepository(db)
	authService := auth.NewService(authRepo, cfg.SMTP)
	authHandler := handler.NewAuthHandler(authService)
	authRoutes := api.Group("/auth")
	authRoutes.Post("/signup", publicRateLimiter, authHandler.CreateUser)
	authRoutes.Get("/verify-email", publicRateLimiter, authHandler.VerifyEmail)
	authRoutes.Post("/login", loginRateLimiter, authHandler.LoginUser)
	authRoutes.Post("/forgot-password", publicRateLimiter, authHandler.ForgotPassword)
	authRoutes.Post("/reset-password", publicRateLimiter, authHandler.ResetPassword)
	authRoutes.Post("/refresh", publicRateLimiter, authHandler.RefreshToken)
	authRoutes.Post("/logout", jwt, authHandler.Logout)

	// User APIs
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, db)
	userHandler := handler.NewUserHandler(userService)
	userRoutes := api.Group("/users")
	userRoutes.Get("/", jwt, authRateLimiter, userHandler.GetAllUsers)
	userRoutes.Get("/:id", jwt, authRateLimiter, userHandler.GetUserByID)
	userRoutes.Patch("/:id", jwt, authRateLimiter, userHandler.UpdateUser)
	userRoutes.Delete("/:id", jwt, authRateLimiter, middleware.RBAC(constants.RoleAdmin), userHandler.DeleteUser)

	// Public routes
	publicHandler := handler.NewPublicHandler(db)
	api.Get("/health/server", publicRateLimiter, publicHandler.HealthCheckServer)
	api.Get("/health/db", publicRateLimiter, publicHandler.HealthCheckDB)
}
