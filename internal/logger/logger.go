package logger

import (
	"os"
	"strings"
	"time"

	"github.com/sudo-hassan-zahid/go-api-server/internal/config"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Log = log.Logger

func Init(cfg config.LogConfig, env string) {
	level, err := zerolog.ParseLevel(strings.ToLower(cfg.Level))
	if err != nil {
		level = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(level)

	if env == "local" {
		Log = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		})
	}

	Log.Info().Str("env", env).Msg("Logger initialized")
}
