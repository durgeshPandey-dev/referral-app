package service

import (
	"context"

	"github.com/durgeshPandey-dev/referral/apps/api/internal/excel"
	"github.com/durgeshPandey-dev/referral/apps/api/internal/queue"
	"github.com/durgeshPandey-dev/referral/apps/api/logger"
)

type referral struct {
	queue *queue.Queue
}

func newReferral(q *queue.Queue) *referral {
	return &referral{queue: q}
}

func (s *referral) ProcessReferral(ctx context.Context, filePath string) error {
	logger.Info(ctx, "process_started", map[string]any{
		"file_path": filePath,
	})

	contacts, err := excel.Parse(filePath)
	if err != nil {
		logger.Error(ctx, "excel_parse_failed", map[string]any{
			"error":     err,
			"file_path": filePath,
		})
		return err
	}

	logger.Info(ctx, "excel_parsed", map[string]any{
		"count": len(contacts),
	})

	for _, c := range contacts {
		select {
		case <-ctx.Done():
			logger.Warn(ctx, "process_cancelled", map[string]any{
				"reason": ctx.Err(),
			})
			return ctx.Err()

		default:
			s.queue.Enqueue(queue.Job{
				Contact: c,
				Ctx:     ctx,
			})

			logger.Info(ctx, "job_enqueued", map[string]any{
				"email": c.Email,
			})
		}
	}

	logger.Info(ctx, "process_completed", map[string]any{
		"total_jobs": len(contacts),
	})

	return nil
}
