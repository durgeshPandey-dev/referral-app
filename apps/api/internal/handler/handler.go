package handler

import (
	v1 "github.com/durgeshPandey-dev/referral/apps/api/internal/handler/v1"
	"github.com/durgeshPandey-dev/referral/apps/api/internal/service"
	"github.com/gofiber/fiber/v2"
)

type Version int

const (
	nilVersion Version = iota
	V1
)

type Registrar interface {
	Register(fiber.Router)
}

type handler struct {
	registrar
}

type registrar = map[Version]Registrar

func New(svc *service.Service) *handler {
	hls := make(registrar)
	hls[V1] = v1.New(svc)

	return &handler{hls}
}

func (h *handler) Register(app *fiber.App) {
	app.Get("/health", func(c *fiber.Ctx) error {
		_, err := c.WriteString("OK")
		return err
	})

	apiGrp := app.Group("/api")
	for _, r := range h.registrar {
		r.Register(apiGrp)
	}
}
