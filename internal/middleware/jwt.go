package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/sudo-hassan-zahid/go-api-server/internal/auth"
	appErrors "github.com/sudo-hassan-zahid/go-api-server/internal/errors"
)

func JWT(jwtAuth *auth.JWT) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return appErrors.HandleError(c, appErrors.ErrUnauthorized)
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return appErrors.HandleError(c, appErrors.ErrUnauthorized)
		}

		claims, err := jwtAuth.Validate(parts[1])
		if err != nil {
			return appErrors.HandleError(c, err)
		}

		c.Locals("userID", claims.UserID)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}
