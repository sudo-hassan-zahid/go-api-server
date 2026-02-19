package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
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

	if strings.ToLower(env) == constants.ENV_DEVELOPMENT {
		console := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "15:04:05",
		}

		console.FormatCaller = func(i interface{}) string {
			if c, ok := i.(string); ok {
				return filepath.Base(c)
			}
			return ""
		}

		console.FormatLevel = func(i interface{}) string {
			level := strings.ToUpper(i.(string))
			switch level {
			case "INFO":
				return "\033[32m" + level + "\033[0m"
			case "WARN":
				return "\033[33m" + level + "\033[0m"
			case "ERROR":
				return "\033[31m" + level + "\033[0m"
			case "DEBUG":
				return "\033[36m" + level + "\033[0m"
			default:
				return level
			}
		}

		console.FormatMessage = func(i interface{}) string {
			return i.(string)
		}

		console.FormatFieldName = func(i interface{}) string {
			return i.(string) + "="
		}

		console.FormatFieldValue = func(i interface{}) string {
			return strings.TrimSpace(fmt.Sprintf("%v", i))
		}

		output = console
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
