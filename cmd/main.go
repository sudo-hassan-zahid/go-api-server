package main

import (
	"context"
	"database/sql"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	_ "github.com/sudo-hassan-zahid/go-api-server/docs"
	"github.com/sudo-hassan-zahid/go-api-server/internal/auth"
	"github.com/sudo-hassan-zahid/go-api-server/internal/config"
	"github.com/sudo-hassan-zahid/go-api-server/internal/constants"
	"github.com/sudo-hassan-zahid/go-api-server/internal/database"
	"github.com/sudo-hassan-zahid/go-api-server/internal/handler"
	appLogger "github.com/sudo-hassan-zahid/go-api-server/internal/logger"
	"github.com/sudo-hassan-zahid/go-api-server/internal/models"
	"github.com/sudo-hassan-zahid/go-api-server/routes"
	swagger "github.com/swaggo/fiber-swagger"
)

// @title             				Go API Server
// @version           				1.0
// @description       				This API server is powered by Go. Using PostgreSQL for DB with a magical touch of GORM
// @BasePath          				/api
// @securityDefinitions.apikey		Bearer
// @in 								header
// @name 							Authorization
// @description 					Type "Bearer" followed by your JWT token.
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
	appLogger.Init(cfg.Log, cfg.App.Environment)

	// Initialize Database
	db, err := database.Connect(cfg.DB, cfg.App.Environment == constants.ENV_DEVELOPMENT)
	if err != nil {
		return err
	}

	// Initialize Database
	if cfg.App.Environment == constants.ENV_DEVELOPMENT {
		if err := db.AutoMigrate(&models.User{}); err != nil {
			return err
		}
	}

	// Initialize Fiber App
	app := handler.NewApp()

	// Middlewares
	app.Use(recover.New())
	app.Use(logger.New())

	// Initialize auth
	auth.Init(cfg)

	// Route setup
	routes.Setup(app, db)

	// Swagger
	app.Get("/swagger/*", swagger.FiberWrapHandler())

	// Server errors
	serverErrors := make(chan error, 1)

	go func() {
		appLogger.Log.Info().
			Str("port", cfg.App.Port).
			Msg("Starting server")

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

	appLogger.Log.Info().Msg("Server gracefully stopped")

	return nil
}

func closeDatabase(dbConn interface{ DB() (*sql.DB, error) }) error {
	sqlDB, err := dbConn.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
