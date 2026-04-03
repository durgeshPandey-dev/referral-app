package queue

import (
	"github.com/durgeshPandey-dev/referral/apps/api/internal/email"
	"github.com/durgeshPandey-dev/referral/apps/api/logger"
)

func (q *Queue) StartWorkers(n int) {
	for i := 0; i < n; i++ {
		workerID := i + 1

		go func(id int) {
			for job := range q.Jobs {
				ctx := job.Ctx

				select {
				case <-ctx.Done():
					logger.Warn(ctx, "job_cancelled", map[string]interface{}{
						"email":     job.Contact.Email,
						"worker_id": id,
					})
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
					logger.Error(ctx, "worker_email_failed", map[string]interface{}{
						"email":     job.Contact.Email,
						"worker_id": id,
					})
					continue
				}

				logger.Info(ctx, "worker_email_sent", map[string]interface{}{
					"email":     job.Contact.Email,
					"worker_id": id,
				})
			}
		}(workerID)
	}
}
