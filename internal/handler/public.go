package handler

import "github.com/gofiber/fiber/v2"

type PublicHandler struct{}

func NewPublicHandler() *PublicHandler {
	return &PublicHandler{}
}

// @Summary     Health check
// @Description Checks if the server is up and running
// @Tags        Health
// @Accept      json
// @Produce     json
// @Success 	200 {object} map[string]interface{}
// @Router      /health [get]
func (h *PublicHandler) HealthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "ok",
	})
}
