package queue

import (
	"referral-app/internal/email"
	"referral-app/internal/observability"
	"referral-app/pkg/logger"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func (q *Queue) StartWorkers(n int) {
	tracer := otel.Tracer("beckn-onix/worker")

	for i := 0; i < n; i++ {
		workerID := i + 1

		go func(id int) {
			for job := range q.Jobs {
				ctx := job.Ctx
				ctx, span := tracer.Start(ctx, "worker.process_job")
				span.SetAttributes(
					attribute.Int("worker.id", id),
					attribute.String("job.email", job.Contact.Email),
				)
				observability.SetQueueDepth(len(q.Jobs))

				select {
				case <-ctx.Done():
					observability.IncQueueJobEvent("cancelled")
					logger.Warn(ctx, "job_cancelled", map[string]interface{}{
						"email":     job.Contact.Email,
						"worker_id": id,
					})
					span.SetStatus(codes.Error, "job cancelled")
					span.End()
					continue
				default:
				}

				logger.Info(ctx, "worker_started_job", map[string]interface{}{
					"email":     job.Contact.Email,
					"worker_id": id,
				})

				content := email.BuildTemplate(
					job.Contact.Name,
					job.Contact.CompanyName,
				)

				err := email.Send(ctx, job.Contact.Email, "Opportunity", content)
				if err != nil {
					observability.IncQueueJobEvent("failed")
					logger.Error(ctx, "worker_email_failed", map[string]interface{}{
						"email":     job.Contact.Email,
						"worker_id": id,
					})
					span.RecordError(err)
					span.SetStatus(codes.Error, err.Error())
					span.End()
					continue
				}

				observability.IncQueueJobEvent("sent")
				logger.Info(ctx, "worker_email_sent", map[string]interface{}{
					"email":     job.Contact.Email,
					"worker_id": id,
				})
				span.SetStatus(codes.Ok, "email sent")
				span.End()
			}
		}(workerID)
	}
}
