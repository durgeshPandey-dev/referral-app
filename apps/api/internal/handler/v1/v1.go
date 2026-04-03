package v1

import (
	"context"
	"path/filepath"

	"github.com/durgeshPandey-dev/referral/apps/api/internal/service"
	"github.com/durgeshPandey-dev/referral/apps/api/internal/utils"
	"github.com/durgeshPandey-dev/referral/apps/api/logger"
	"github.com/gofiber/fiber/v2"
)

type v1 struct {
	svc *service.Service
}

func New(svc *service.Service) *v1 {
	return &v1{svc: svc}
}

func (v1 *v1) Register(apiGrp fiber.Router) {
	v1Grp := apiGrp.Group("/v1")

	v1Grp.Post("/upload", v1.upload)
}

// Upload godoc
// @Summary      Upload an Excel containing contact details
// @Accept       multipart/form-data
// @Produce      json
// @Param        file  formData  file  true  "Excel (.xlsx) file with contact details"
// @Success      200   {object}  map[string]string
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /v1/upload [post]
func (v *v1) upload(c *fiber.Ctx) error {
	ctxVal := c.Locals("ctx")
	ctx, ok := ctxVal.(context.Context)
	if !ok {
		ctx = context.Background()
	}

	file, err := c.FormFile("file")
	if err != nil {
		logger.Error(ctx, "file_missing", map[string]any{
			"error": err,
		})
		return c.Status(400).JSON(fiber.Map{"error": "file required"})
	}

	if filepath.Ext(file.Filename) != ".xlsx" {
		logger.Warn(ctx, "invalid_file_type", map[string]any{
			"filename": file.Filename,
		})
		return c.Status(400).JSON(fiber.Map{"error": "only .xlsx allowed"})
	}

	path, err := utils.SaveUploadedFile(file, "./uploads")
	if err != nil {
		logger.Error(ctx, "file_save_failed", map[string]any{
			"error": err,
		})
		return c.Status(500).JSON(fiber.Map{"error": "failed to save"})
	}

	logger.Info(ctx, "file_uploaded", map[string]any{
		"file": file.Filename,
		"path": path,
	})

	go func(parentCtx context.Context) {
		// extract request id
		reqID := parentCtx.Value(logger.RequestIDKey)

		// new background ctx
		bgCtx := context.Background()

		// reattach request id
		if reqID != nil {
			bgCtx = context.WithValue(bgCtx, logger.RequestIDKey, reqID)
		}

		err := v.svc.ProcessReferral(bgCtx, path)
		if err != nil {
			logger.Error(bgCtx, "background_process_failed", map[string]any{
				"error": err,
			})
		}
	}(ctx)

	return c.JSON(fiber.Map{
		"message": "processing started",
	})
}
