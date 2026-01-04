package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	_ "github.com/sudo-hassan-zahid/go-api-server/docs"
	"github.com/sudo-hassan-zahid/go-api-server/internal/auth"
	"github.com/sudo-hassan-zahid/go-api-server/internal/config"
	"github.com/sudo-hassan-zahid/go-api-server/internal/constants"
	"github.com/sudo-hassan-zahid/go-api-server/internal/database"
	appLogger "github.com/sudo-hassan-zahid/go-api-server/internal/logger"
	"github.com/sudo-hassan-zahid/go-api-server/internal/middleware"
	"github.com/sudo-hassan-zahid/go-api-server/internal/models"
	"github.com/sudo-hassan-zahid/go-api-server/routes"
	swagger "github.com/swaggo/fiber-swagger"
)

// @title            				Go API Server
// @version          				1.0
// @description      				This API server is powered by Go. Using PostgreSQL for DB with a magical touch of GORM
// @BasePath         				/api
// @securityDefinitions.apikey  	BearerAuth
// @in 								header
// @name 							Authorization
// @description 					Type "Bearer" followed by your JWT token.
func main() {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	// Initialize Logger
	appLogger.Init(cfg.Log, cfg.App.Environment)

	// Connect to Database
	db, err := database.Connect(cfg.DB, cfg.App.Environment == constants.ENV_LOCAL)
	if err != nil {
		appLogger.Log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	// Auto-migrate dev models
	if cfg.App.Environment == constants.ENV_LOCAL {
		if err := db.AutoMigrate(&models.User{}); err != nil {
			log.Fatal("AutoMigrate failed:", err)
		}
	}

	// Initialize Fiber App
	app := fiber.New(fiber.Config{
		AppName:      cfg.App.Name,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
	})

	// Middlewares
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(middleware.ErrorLogger())

	// Auth init
	auth.Init(cfg)

	// Routes
	routes.Setup(app, db)

	// Swagger docs
	app.Get("/swagger/*", swagger.FiberWrapHandler())

	// Start Server in Goroutine
	serverErrors := make(chan error, 1)
	go func() {
		appLogger.Log.Info().Str("port", cfg.App.Port).Msg("Starting Fiber server")
		serverErrors <- app.Listen(":" + cfg.App.Port)
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-quit:
		appLogger.Log.Info().Str("signal", sig.String()).Msg("Shutting down server...")
	case err := <-serverErrors:
		appLogger.Log.Fatal().Err(err).Msg("Server failed")
	}

	// Fiber shutdown with timeout context
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.Shutdown(); err != nil {
		appLogger.Log.Error().Err(err).Msg("Error during server shutdown")
	} else {
		appLogger.Log.Info().Msg("Server gracefully stopped")
	}

	// Close DB connection
	sqlDB, _ := db.DB()
	if err := sqlDB.Close(); err != nil {
		appLogger.Log.Error().Err(err).Msg("Failed to close database connection")
	} else {
		appLogger.Log.Info().Msg("Database connection closed")
	}
}
