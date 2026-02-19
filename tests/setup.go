package tests

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sudo-hassan-zahid/go-api-server/internal/auth"
	"github.com/sudo-hassan-zahid/go-api-server/internal/config"
	"github.com/sudo-hassan-zahid/go-api-server/internal/database"
	"github.com/sudo-hassan-zahid/go-api-server/routes"
	"github.com/testcontainers/testcontainers-go"
	tcPostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	tcRedis "github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
	driverPostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	TestDB         *gorm.DB
	TestRdb        *redis.Client
	TestApp        *fiber.App
	PGContainer    *tcPostgres.PostgresContainer
	RedisContainer *tcRedis.RedisContainer
)

func SetupTestContainer() (func(), error) {
	ctx := context.Background()

	// Redis Container
	redisC, err := tcRedis.Run(ctx, "redis:7-alpine")
	if err != nil {
		return nil, fmt.Errorf("failed to start redis container: %w", err)
	}
	RedisContainer = redisC

	redisHost, err := redisC.Host(ctx)
	if err != nil {
		return nil, err
	}
	redisPort, err := redisC.MappedPort(ctx, "6379")
	if err != nil {
		return nil, err
	}

	// Connect to Redis
	TestRdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", redisHost, redisPort.Port()),
	})

	// Set the global Redis client so the app uses this one
	database.Rdb = TestRdb

	// Postgres Container
	pgC, err := tcPostgres.Run(ctx, "postgres:16-alpine",
		tcPostgres.WithDatabase("testdb"),
		tcPostgres.WithUsername("testuser"),
		tcPostgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start postgres container: %w", err)
	}
	PGContainer = pgC

	pgHost, err := pgC.Host(ctx)
	if err != nil {
		return nil, err
	}
	pgPort, err := pgC.MappedPort(ctx, "5432")
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("host=%s user=testuser password=testpass dbname=testdb port=%s sslmode=disable TimeZone=UTC", pgHost, pgPort.Port())
	TestDB, err = gorm.Open(driverPostgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to test db: %w", err)
	}

	cleanup := func() {
		if err := redisC.Terminate(ctx); err != nil {
			log.Printf("failed to terminate redis container: %s", err)
		}
		if err := pgC.Terminate(ctx); err != nil {
			log.Printf("failed to terminate postgres container: %s", err)
		}
	}

	return cleanup, nil
}

func SetupApp() *fiber.App {
	app := fiber.New()

	cfg := &config.Config{
		App: config.AppConfig{
			Environment: "test",
			JWTSecret:   []byte("secret"),
		},
		DB:    config.DBConfig{},
		Redis: config.RedisConfig{},
		SMTP:  config.SMTPConfig{},
	}

	auth.Init(cfg)

	routes.Setup(app, TestDB, cfg)

	return app
}
