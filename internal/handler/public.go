package handler

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type PublicHandler struct {
	db *gorm.DB
}

func NewPublicHandler(db *gorm.DB) *PublicHandler {
	return &PublicHandler{db: db}
}

// @Summary     Health check server
// @Description Checks if the server is up and running
// @Tags        Health
// @Accept      json
// @Produce     json
// @Success 	200 {object} map[string]interface{}
// @Success 	429 {object} map[string]string "Too many requests"
// @Router      /health/server [get]
func (h *PublicHandler) HealthCheckServer(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "ok",
		"service": "api",
		"time":    time.Now().UTC(),
	})
}

// @Summary     Health check database
// @Description Checks if the database is up and running
// @Tags        Health
// @Accept      json
// @Produce     json
// @Success     200 {object} map[string]interface{}
// @Failure     429 {object} map[string]string "Too many requests"
// @Failure     503 {object} map[string]interface{}
// @Router      /health/db [get]
func (h *PublicHandler) HealthCheckDB(c *fiber.Ctx) error {
	sqlDB, err := h.db.DB()
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status":    "error",
			"component": "database",
			"error":     err.Error(),
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status":    "error",
			"component": "database",
			"error":     "database unreachable",
		})
	}

	stats := sqlDB.Stats()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "ok",
		"database": fiber.Map{
			"open_connections": stats.OpenConnections,
			"in_use":           stats.InUse,
			"idle":             stats.Idle,
		},
	})
}
