package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/gorm"

	_ "github.com/sudo-hassan-zahid/go-api-server/docs"
	"github.com/sudo-hassan-zahid/go-api-server/internal/auth"
	"github.com/sudo-hassan-zahid/go-api-server/internal/config"
	"github.com/sudo-hassan-zahid/go-api-server/internal/constants"
	"github.com/sudo-hassan-zahid/go-api-server/internal/database"
	"github.com/sudo-hassan-zahid/go-api-server/internal/handler"
	appLogger "github.com/sudo-hassan-zahid/go-api-server/internal/logger"
	"github.com/sudo-hassan-zahid/go-api-server/internal/middleware"
	"github.com/sudo-hassan-zahid/go-api-server/routes"
	swagger "github.com/swaggo/fiber-swagger"
)

// @title             				Go API Server
// @version           				1.0
// @description       				This is a production-grade API server written in Go using Fiber and GORM.
// @termsOfService    				http://swagger.io/terms/
// @license.name      				Apache 2.0
// @license.url       				http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath          				/api
// @securityDefinitions.apikey		Bearer
// @in 								header
// @name 							Authorization
// @description 					Type "Bearer" followed by a space and JWT token.
func main() {
	if err := run(); err != nil {
		appLogger.Log.Fatal().Err(err).Msg("Application crashed")
	}
}

func run() error {
	// Load Config
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	// Initialize Logger
	appLogger.Init(cfg.Log.Level, cfg.App.Environment)

	// Initialize Database
	db, err := database.Connect(cfg.DB, cfg.App.Environment == constants.ENV_DEVELOPMENT)
	if err != nil {
		return err
	}

	// Initialize Redis
	if err := database.ConnectRedis(cfg.Redis); err != nil {
		return err
	}

	// Initialize Database
	if cfg.App.Environment == constants.ENV_DEVELOPMENT {
		shouldSeed := os.Getenv("SEED_DB") == "true"
		if err := database.Setup(db, shouldSeed); err != nil {
			return err
		}
	}

	// Initialize Fiber App
	app := handler.NewApp()

	// Middlewares
	app.Use(middleware.RequestLogger())
	app.Use(recover.New(recover.Config{
		EnableStackTrace: cfg.App.Environment == constants.ENV_DEVELOPMENT,
	}))

	// Initialize auth
	auth.Init(cfg)

	// Route setup
	routes.Setup(app, db, cfg)

	// Swagger
	app.Get("/swagger/*", swagger.FiberWrapHandler())

	// Server errors
	serverErrors := make(chan error, 1)

	go func() {
		appLogger.Log.Info().
			Str("port", cfg.App.Port).
			Str("env", cfg.App.Environment).
			Msg("starting server")

		serverErrors <- app.Listen(":" + cfg.App.Port)
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-quit:
		appLogger.Log.Info().
			Str("signal", sig.String()).
			Msg("Shutdown signal received")

	case err := <-serverErrors:
		return err
	}

	// Shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		appLogger.Log.Error().Err(err).Msg("Server shutdown error")
	}

	// Close Database
	if err := closeDatabase(db); err != nil {
		appLogger.Log.Error().Err(err).Msg("Failed to close DB")
	}

	// Close Redis
	if err := database.CloseRedis(); err != nil {
		appLogger.Log.Error().Err(err).Msg("Failed to close Redis")
	}

	appLogger.Log.Info().Msg("Server gracefully stopped")

	return nil
}

func closeDatabase(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
