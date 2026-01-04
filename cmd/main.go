package main

import (
	"os"
	"os/signal"
	"syscall"

	stdlog "log"

	"github.com/rs/zerolog/log"

	"github.com/sudo-hassan-zahid/go-api-server/internal/config"
	"github.com/sudo-hassan-zahid/go-api-server/internal/constant"
	"github.com/sudo-hassan-zahid/go-api-server/internal/database"
	appLogger "github.com/sudo-hassan-zahid/go-api-server/internal/logger"
	"github.com/sudo-hassan-zahid/go-api-server/internal/models"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		stdlog.Fatal(err)
	}

	appLogger.Init(cfg.Log, cfg.App.Environment)

	db, err := database.Connect(cfg.DB, cfg.App.Environment == constant.ENV_DEV)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	if cfg.App.Environment == constant.ENV_DEV {
		if err := db.AutoMigrate(&models.User{}); err != nil {
			log.Fatal().Err(err).Msg("database migration failed")
		}
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get sql.DB")
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("shutting down application")

	if err := sqlDB.Close(); err != nil {
		log.Error().Err(err).Msg("failed to close database")
	}

	log.Info().Msg("shutdown complete")
}
