package main

import (
	"fmt"
	"time"

	"github.com/durgeshPandey-dev/referral/apps/api/internal/config"
	"github.com/durgeshPandey-dev/referral/apps/api/internal/handler"
	"github.com/durgeshPandey-dev/referral/apps/api/internal/middleware"
	"github.com/durgeshPandey-dev/referral/apps/api/internal/queue"
	"github.com/durgeshPandey-dev/referral/apps/api/internal/service"

	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
)

// @title           Swagger: Referral Backend API
// @version         1.0
// @description     This is the backend API for the Referral Email Automation System.

// @contact.name   Durgesh Pandey
// @contact.url    https://github.com/durgeshPandey-dev/referral-app

// @host      localhost:3000
// @BasePath  /api
func main() {
	cfg := config.Load()

	app := fiber.New()

	// Middlewares
	app.Use(middleware.LoggerMiddleware(
		time.Duration(cfg.RequestTimeoutSeconds) * time.Second,
	))

	app.Use(swagger.New(swagger.Config{
		FilePath: "../../docs/api/v1/swagger.json",
		Path:     "docs/v1",
		Title:    "Referral Backend API Documentation",
		CacheAge: 1800,
	}))

	// Queue + Workers
	q := queue.NewQueue(cfg.QueueSize)
	q.StartWorkers(cfg.WorkerCount)

	// Services
	svc := service.New(q)

	// Handlers
	h := handler.New(svc)
	h.Register(app)

	fmt.Println("Server running on port:", cfg.Port)

	if err := app.Listen(":" + cfg.Port); err != nil {
		panic(err)
	}
}
