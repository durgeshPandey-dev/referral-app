package queue

import (
	"context"

	"referral-app/internal/models"
	"referral-app/internal/observability"
)

type Job struct {
	Contact models.HRContact
	Ctx     context.Context
}

type Queue struct {
	Jobs chan Job
}

func NewQueue(size int) *Queue {
	q := &Queue{
		Jobs: make(chan Job, size),
	}
	observability.SetQueueDepth(0)
	return q
}

func (q *Queue) Enqueue(job Job) {
	q.Jobs <- job
	observability.IncQueueJobEvent("enqueued")
	observability.SetQueueDepth(len(q.Jobs))
}
