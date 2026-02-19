package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sudo-hassan-zahid/go-api-server/internal/constants"
)

func RBAC(allowedRoles ...constants.Role) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole, ok := c.Locals("role").(string)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Access denied: No role found",
			})
		}

		for _, role := range allowedRoles {
			if string(role) == userRole {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied: Insufficient permissions",
		})
	}
}
