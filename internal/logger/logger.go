package logger

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/sudo-hassan-zahid/go-api-server/internal/constants"
)

var Log zerolog.Logger

func Init(levelStr string, env string) {
	level, err := zerolog.ParseLevel(strings.ToLower(levelStr))
	if err != nil {
		level = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(level)
	zerolog.TimeFieldFormat = time.RFC3339

	var output io.Writer = os.Stdout

	if env == constants.ENV_DEVELOPMENT {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
	}

	Log = zerolog.New(output).
		With().
		Timestamp().
		Caller().
		Str("service", "go-api-server").
		Str("env", env).
		Logger()

	Log.Info().Msg("logger initialized")
}
