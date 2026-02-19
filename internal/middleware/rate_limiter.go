package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/sudo-hassan-zahid/go-api-server/internal/database"
)

func PublicRateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        50,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"error":   "too many requests, please try again later",
			})
		},
	})
}

func AuthRateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			userID := c.Locals("userID")
			if userID != nil {
				return userID.(string)
			}
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"error":   "rate limit exceeded",
			})
		},
	})
}

func LoginRateLimiter() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ip := c.IP()
		key := "login_attempts:" + ip

		limit := 5
		duration := 15 * time.Minute

		count, err := database.Rdb.Incr(c.Context(), key).Result()
		if err != nil {
			return c.Next()
		}

		if count == 1 {
			database.Rdb.Expire(c.Context(), key, duration)
		}

		if count > int64(limit) {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"error":   "Too many login attempts. Please try again later.",
			})
		}

		return c.Next()
	}
}
