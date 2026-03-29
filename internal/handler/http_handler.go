package handler

import (
	"context"
	"path/filepath"

	"referral-app/internal/observability"
	"referral-app/internal/service"
	"referral-app/internal/utils"
	"referral-app/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func RegisterRoutes(app *fiber.App, svc *service.ReferralService) {
	tracer := otel.Tracer("beckn-onix/handler")

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	app.Post("/upload", func(c *fiber.Ctx) error {
		ctxVal := c.Locals("ctx")
		ctx, ok := ctxVal.(context.Context)
		if !ok {
			ctx = context.Background()
		}

		ctx, span := tracer.Start(ctx, "upload.request")
		defer span.End()
		c.Locals("ctx", ctx)

		file, err := c.FormFile("file")
		if err != nil {
			observability.IncUploadRequest("bad_request")
			span.RecordError(err)
			span.SetStatus(codes.Error, "missing file")
			logger.Error(ctx, "file_missing", map[string]interface{}{
				"error": err,
			})
			return c.Status(400).JSON(fiber.Map{"error": "file required"})
		}

		if filepath.Ext(file.Filename) != ".xlsx" {
			observability.IncUploadRequest("bad_request")
			span.SetStatus(codes.Error, "invalid file extension")
			logger.Warn(ctx, "invalid_file_type", map[string]interface{}{
				"filename": file.Filename,
			})
			return c.Status(400).JSON(fiber.Map{"error": "only .xlsx allowed"})
		}

		// ✅ USE UTILS HERE
		path, err := utils.SaveUploadedFile(file, "./uploads")
		if err != nil {
			observability.IncUploadRequest("failed")
			span.RecordError(err)
			span.SetStatus(codes.Error, "file save failed")
			logger.Error(ctx, "file_save_failed", map[string]interface{}{
				"error": err,
			})
			return c.Status(500).JSON(fiber.Map{"error": "failed to save"})
		}

		span.SetAttributes(
			attribute.String("upload.filename", file.Filename),
			attribute.String("upload.path", path),
		)
		observability.IncUploadRequest("accepted")
		span.SetStatus(codes.Ok, "accepted")

		logger.Info(ctx, "file_uploaded", map[string]interface{}{
			"file": file.Filename,
			"path": path,
		})

		go func(parentCtx context.Context) {
			bgCtx := context.WithoutCancel(parentCtx)
			if reqID := parentCtx.Value(logger.RequestIDKey); reqID != nil {
				bgCtx = context.WithValue(bgCtx, logger.RequestIDKey, reqID)
			}

			err := svc.Process(bgCtx, path)
			if err != nil {
				logger.Error(bgCtx, "background_process_failed", map[string]interface{}{
					"error": err,
				})
			}
		}(ctx)

		return c.JSON(fiber.Map{
			"message": "processing started",
		})
	})
}
