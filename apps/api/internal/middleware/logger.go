package middleware

import (
	"context"
	"time"

	"github.com/durgeshPandey-dev/referral/apps/api/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func LoggerMiddleware(timeout time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		reqID := uuid.NewString()

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		ctx = context.WithValue(ctx, logger.RequestIDKey, reqID)
		c.Locals("ctx", ctx)

		start := time.Now()

		err := c.Next()

		logger.Info(ctx, "request_complete", map[string]any{
			"path":     c.Path(),
			"method":   c.Method(),
			"duration": time.Since(start).Milliseconds(),
			"status":   c.Response().StatusCode(),
		})

		return err
	}
}
