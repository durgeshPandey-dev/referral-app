package service

import "github.com/durgeshPandey-dev/referral/apps/api/internal/queue"

type Service struct {
	*services
}

type services struct {
	*referral
}

func New(queue *queue.Queue) *Service {
	svc := &Service{
		&services{
			referral: newReferral(queue),
		},
	}

	return svc
}
