package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"referral-app/internal/observability"
	"referral-app/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
)

type fiberRequestHeaderCarrier struct {
	c *fiber.Ctx
}

func (f fiberRequestHeaderCarrier) Get(key string) string {
	return f.c.Get(key)
}

func (f fiberRequestHeaderCarrier) Set(key, value string) {
	// The inbound carrier exists only for extraction; setter is required by interface.
}

func (f fiberRequestHeaderCarrier) Keys() []string {
	headers := f.c.GetReqHeaders()
	keys := make([]string, 0, len(headers))
	for k := range headers {
		keys = append(keys, strings.TrimSpace(k))
	}
	return keys
}

func RequestMiddleware(timeout time.Duration) fiber.Handler {
	tracer := otel.Tracer("beckn-onix/http")

	return func(c *fiber.Ctx) error {
		// Fiber reuses request buffers; clone strings used after c.Next().
		method := strings.Clone(c.Method())
		requestPath := strings.Clone(c.Path())

		reqID := c.Get("X-Request-Id")
		if reqID == "" {
			reqID = uuid.NewString()
		}
		c.Set("X-Request-Id", reqID)

		parentCtx := otel.GetTextMapPropagator().Extract(
			context.Background(),
			propagation.TextMapCarrier(fiberRequestHeaderCarrier{c: c}),
		)
		ctx, cancel := context.WithTimeout(parentCtx, timeout)
		defer cancel()

		ctx = context.WithValue(ctx, logger.RequestIDKey, reqID)
		c.Locals("ctx", ctx)

		start := time.Now()
		routeLabel := observability.NormalizeRouteLabel(requestPath)
		observability.IncInflightHTTPRequests()

		ctx, span := tracer.Start(ctx, "http.request")
		span.SetAttributes(
			attribute.String("http.method", method),
			attribute.String("http.route", routeLabel),
			attribute.String("request.id", reqID),
		)
		c.Locals("ctx", ctx)

		defer func() {
			duration := time.Since(start)
			status := c.Response().StatusCode()
			observability.DecInflightHTTPRequests()
			observability.ObserveHTTPRequest(method, routeLabel, status, duration)

			span.SetAttributes(
				attribute.Int("http.status_code", status),
				attribute.Int64("http.duration_ms", duration.Milliseconds()),
			)
			if status >= http.StatusBadRequest {
				span.SetStatus(codes.Error, http.StatusText(status))
			} else {
				span.SetStatus(codes.Ok, "success")
			}
			span.End()
		}()

		err := c.Next()
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}

		logger.Info(ctx, "request_complete", map[string]interface{}{
			"path":     requestPath,
			"method":   method,
			"duration": time.Since(start).Milliseconds(),
			"status":   c.Response().StatusCode(),
		})

		return err
	}
}
