package queue

import (
	"context"

	"github.com/durgeshPandey-dev/referral/apps/api/internal/models"
)

type Job struct {
	Contact models.HRContact
	Ctx     context.Context
}

type Queue struct {
	Jobs chan Job
}

func NewQueue(size int) *Queue {
	return &Queue{
		Jobs: make(chan Job, size),
	}
}

func (q *Queue) Enqueue(job Job) {
	q.Jobs <- job
}
