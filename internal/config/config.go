package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/sudo-hassan-zahid/go-api-server/internal/constants"
)

type AppConfig struct {
	Name        string
	Environment string
	Port        string
	JWTSecret   []byte
}

type DBConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type LogConfig struct {
	Level string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type SMTPConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	From     string
}

type Config struct {
	App   AppConfig
	DB    DBConfig
	Redis RedisConfig
	Log   LogConfig
	SMTP  SMTPConfig
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	jwtSecret, err := requiredEnv("JWT_SECRET")
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		App: AppConfig{
			Name:        getEnv("APP_NAME", "go_api_server"),
			Environment: getEnv("APP_ENVIRONMENT", constants.ENV_DEVELOPMENT),
			Port:        getEnv("APP_PORT", "8080"),
			JWTSecret:   []byte(jwtSecret),
		},
		DB: DBConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "5432"),
			User:            getEnv("DB_USER", "postgres"),
			Password:        getEnv("DB_PASSWORD", ""),
			Name:            getEnv("DB_NAME", "go_api_server"),
			SSLMode:         getEnv("DB_SSLMODE", "disable"),
			MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 25, nil),
			MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 25, nil),
			ConnMaxLifetime: getEnvAsDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute, nil),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0, nil),
		},
		Log: LogConfig{
			Level: getEnv("LOG_LEVEL", "debug"),
		},
		SMTP: SMTPConfig{
			Host:     getEnv("SMTP_HOST", "smtp.mailtrap.io"),
			Port:     getEnvAsInt("SMTP_PORT", 2525, nil),
			User:     getEnv("SMTP_USER", ""),
			Password: getEnv("SMTP_PASSWORD", ""),
			From:     getEnv("SMTP_FROM", "noreply@example.com"),
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func requiredEnv(key string) (string, error) {
	if val := strings.TrimSpace(os.Getenv(key)); val != "" {
		return val, nil
	}
	return "", fmt.Errorf("missing required env: %s", key)
}

func getEnvAsInt(key string, defaultVal int, errs *[]error) int {
	if val := os.Getenv(key); val != "" {
		i, err := strconv.Atoi(val)
		if err != nil {
			if errs != nil {
				*errs = append(*errs, fmt.Errorf("%s must be an integer", key))
				return defaultVal
			}
			return defaultVal
		}
		return i
	}
	return defaultVal
}

func getEnvAsDuration(key string, defaultVal time.Duration, errs *[]error) time.Duration {
	if val := os.Getenv(key); val != "" {
		d, err := time.ParseDuration(val)
		if err != nil {
			if errs != nil {
				*errs = append(*errs, fmt.Errorf("%s must be a duration", key))
				return defaultVal
			}
			return defaultVal
		}
		return d
	}
	return defaultVal
}

func (cfg *Config) validate() error {
	var errs []error

	cfg.DB.MaxOpenConns = getEnvAsInt("DB_MAX_OPEN_CONNS", cfg.DB.MaxOpenConns, &errs)
	cfg.DB.MaxIdleConns = getEnvAsInt("DB_MAX_IDLE_CONNS", cfg.DB.MaxIdleConns, &errs)
	cfg.DB.ConnMaxLifetime = getEnvAsDuration("DB_CONN_MAX_LIFETIME", cfg.DB.ConnMaxLifetime, &errs)
	cfg.Redis.DB = getEnvAsInt("REDIS_DB", cfg.Redis.DB, &errs)
	cfg.SMTP.Port = getEnvAsInt("SMTP_PORT", cfg.SMTP.Port, &errs)

	if strings.TrimSpace(cfg.App.Port) == "" {
		errs = append(errs, fmt.Errorf("APP_PORT is required"))
	}
	if cfg.App.Environment == constants.ENV_PRODUCTION && len(cfg.App.JWTSecret) < 32 {
		errs = append(errs, fmt.Errorf("JWT_SECRET must be at least 32 bytes in production"))
	}
	if cfg.DB.MaxOpenConns < 1 {
		errs = append(errs, fmt.Errorf("DB_MAX_OPEN_CONNS must be greater than zero"))
	}
	if cfg.DB.MaxIdleConns < 1 {
		errs = append(errs, fmt.Errorf("DB_MAX_IDLE_CONNS must be greater than zero"))
	}
	if cfg.DB.MaxIdleConns > cfg.DB.MaxOpenConns {
		errs = append(errs, fmt.Errorf("DB_MAX_IDLE_CONNS cannot exceed DB_MAX_OPEN_CONNS"))
	}

	if len(errs) > 0 {
		return errorsJoin(errs)
	}
	return nil
}

func errorsJoin(errs []error) error {
	if len(errs) == 1 {
		return errs[0]
	}

	messages := make([]string, 0, len(errs))
	for _, err := range errs {
		messages = append(messages, err.Error())
	}
	return fmt.Errorf("invalid config: %s", strings.Join(messages, "; "))
}
