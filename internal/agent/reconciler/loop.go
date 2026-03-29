package reconciler

import (
	"context"
	"log"
	"time"

	"ospnet/internal/agent/runtime/manager"
)

type Loop struct {
	manager  *manager.Manager
	logger   *log.Logger
	interval time.Duration
}

func New(mgr *manager.Manager, logger *log.Logger) *Loop {
	return &Loop{
		manager:  mgr,
		logger:   logger,
		interval: 12 * time.Second,
	}
}

func (l *Loop) Run(ctx context.Context) error {
	ticker := time.NewTicker(l.interval)
	defer ticker.Stop()

	l.runOnce(ctx)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			l.runOnce(ctx)
		}
	}
}

func (l *Loop) runOnce(ctx context.Context) {
	if err := l.manager.Reconcile(ctx); err != nil {
		l.logger.Printf("reconcile failed: %v", err)
	}
}
