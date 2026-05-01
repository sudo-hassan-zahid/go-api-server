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
	deps := newDependencies(db, cfg)

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	api := app.Group("/api")

	registerAuthRoutes(api, deps)
	registerUserRoutes(api, deps)
	registerPublicRoutes(api, deps)
}

type dependencies struct {
	jwt               fiber.Handler
	authRateLimiter   fiber.Handler
	publicRateLimiter fiber.Handler
	loginRateLimiter  fiber.Handler
	authHandler       *handler.AuthHandler
	userHandler       *handler.UserHandler
	publicHandler     *handler.PublicHandler
}

func newDependencies(db *gorm.DB, cfg *config.Config) dependencies {
	userRepo := repository.NewUserRepository(db)

	return dependencies{
		jwt:               middleware.JWTMiddleware(),
		authRateLimiter:   middleware.AuthRateLimiter(),
		publicRateLimiter: middleware.PublicRateLimiter(),
		loginRateLimiter:  middleware.LoginRateLimiter(),
		authHandler:       handler.NewAuthHandler(auth.NewService(userRepo, cfg.SMTP)),
		userHandler:       handler.NewUserHandler(service.NewUserService(userRepo, db)),
		publicHandler:     handler.NewPublicHandler(db),
	}
}

func registerAuthRoutes(api fiber.Router, deps dependencies) {
	authRoutes := api.Group("/auth")
	authRoutes.Post("/signup", deps.publicRateLimiter, deps.authHandler.CreateUser)
	authRoutes.Get("/verify-email", deps.publicRateLimiter, deps.authHandler.VerifyEmail)
	authRoutes.Post("/login", deps.loginRateLimiter, deps.authHandler.LoginUser)
	authRoutes.Post("/forgot-password", deps.publicRateLimiter, deps.authHandler.ForgotPassword)
	authRoutes.Post("/reset-password", deps.publicRateLimiter, deps.authHandler.ResetPassword)
	authRoutes.Post("/refresh", deps.publicRateLimiter, deps.authHandler.RefreshToken)
	authRoutes.Post("/logout", deps.jwt, deps.authHandler.Logout)
}

func registerUserRoutes(api fiber.Router, deps dependencies) {
	userRoutes := api.Group("/users")
	userRoutes.Get("/", deps.jwt, deps.authRateLimiter, deps.userHandler.GetAllUsers)
	userRoutes.Get("/:id", deps.jwt, deps.authRateLimiter, deps.userHandler.GetUserByID)
	userRoutes.Patch("/:id", deps.jwt, deps.authRateLimiter, deps.userHandler.UpdateUser)
	userRoutes.Delete("/:id", deps.jwt, deps.authRateLimiter, middleware.RBAC(constants.RoleAdmin), deps.userHandler.DeleteUser)
}

func registerPublicRoutes(api fiber.Router, deps dependencies) {
	api.Get("/health/server", deps.publicRateLimiter, deps.publicHandler.HealthCheckServer)
	api.Get("/health/db", deps.publicRateLimiter, deps.publicHandler.HealthCheckDB)
}
