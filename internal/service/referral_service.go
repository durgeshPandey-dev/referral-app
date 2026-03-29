package service

import (
	"context"

	"referral-app/internal/excel"
	"referral-app/internal/observability"
	"referral-app/internal/queue"
	"referral-app/pkg/logger"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type ReferralService struct {
	queue *queue.Queue
}

func NewReferralService(q *queue.Queue) *ReferralService {
	return &ReferralService{queue: q}
}

func (s *ReferralService) Process(ctx context.Context, filePath string) error {
	tracer := otel.Tracer("beckn-onix/service")
	ctx, span := tracer.Start(ctx, "referral.process_upload")
	defer span.End()

	logger.Info(ctx, "process_started", map[string]interface{}{
		"file_path": filePath,
	})

	contacts, err := excel.Parse(filePath)
	if err != nil {
		observability.AddContacts("invalid", 1)
		logger.Error(ctx, "excel_parse_failed", map[string]interface{}{
			"error":     err,
			"file_path": filePath,
		})
		span.RecordError(err)
		span.SetStatus(codes.Error, "excel parse failed")
		return err
	}

	logger.Info(ctx, "excel_parsed", map[string]interface{}{
		"count": len(contacts),
	})
	span.SetAttributes(
		attribute.String("file.path", filePath),
		attribute.Int("contacts.count", len(contacts)),
	)
	observability.AddContacts("parsed", len(contacts))

	for _, c := range contacts {
		jobCtx, jobSpan := tracer.Start(ctx, "referral.enqueue_job")
		jobSpan.SetAttributes(attribute.String("contact.email", c.Email))

		select {
		case <-jobCtx.Done():
			observability.IncQueueJobEvent("cancelled")
			logger.Warn(jobCtx, "process_cancelled", map[string]interface{}{
				"reason": jobCtx.Err(),
			})
			jobSpan.SetStatus(codes.Error, "context cancelled")
			jobSpan.End()
			span.SetStatus(codes.Error, "process cancelled")
			return jobCtx.Err()

		default:
			s.queue.Enqueue(queue.Job{
				Contact: c,
				Ctx:     jobCtx,
			})

			observability.AddContacts("enqueued", 1)
			logger.Info(jobCtx, "job_enqueued", map[string]interface{}{
				"email": c.Email,
			})
			jobSpan.SetStatus(codes.Ok, "job enqueued")
			jobSpan.End()
		}
	}

	logger.Info(ctx, "process_completed", map[string]interface{}{
		"total_jobs": len(contacts),
	})
	span.SetStatus(codes.Ok, "process completed")

	return nil
}
