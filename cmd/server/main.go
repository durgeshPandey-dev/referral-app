package main

import (
	"context"
	"fmt"
	"time"

	"referral-app/internal/config"
	"referral-app/internal/handler"
	"referral-app/internal/middleware"
	"referral-app/internal/observability"
	"referral-app/internal/queue"
	"referral-app/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

func main() {
	cfg := config.Load()

	app := fiber.New()

	if cfg.ObservabilityEnabled {
		observability.InitMetrics()

		shutdownTracing, err := observability.InitTracing(
			context.Background(),
			cfg.ServiceName,
			cfg.Environment,
			cfg.OTLPTraceEndpoint,
			cfg.OTLPInsecure,
		)
		if err != nil {
			panic(err)
		}
		defer func() {
			_ = shutdownTracing(context.Background())
		}()

		app.Get("/metrics", adaptor.HTTPHandler(observability.PrometheusHTTPHandler()))
	}

	// Middleware
	app.Use(middleware.RequestMiddleware(
		time.Duration(cfg.RequestTimeoutSeconds) * time.Second,
	))

	// Queue + Workers
	q := queue.NewQueue(cfg.QueueSize)
	q.StartWorkers(cfg.WorkerCount)

	// Service
	svc := service.NewReferralService(q)

	// Routes
	handler.RegisterRoutes(app, svc)

	fmt.Println("Server running on port:", cfg.Port)

	if err := app.Listen(":" + cfg.Port); err != nil {
		panic(err)
	}
}
