package node

import (
	"context"
	"fmt"

	"ospnet/internal/master/db"
)

type Service interface {
	ListNodes(ctx context.Context) ([]db.Node, error)
	Heartbeat(ctx context.Context, nodeID string) error
}

type sv struct {
	db db.Querier
}

func NewService(db db.Querier) Service {
	return &sv{db: db}
}

func (s *sv) ListNodes(ctx context.Context) ([]db.Node, error) {
	rs, err := s.db.GetNodes(ctx)
	if err != nil {
		return []db.Node{}, fmt.Errorf("Error in listing nodes (%s)", err.Error())
	}

	return rs, nil

}

func (s *sv) Heartbeat(ctx context.Context, nodeID string) error {
	return s.db.UpdateHeartbeat(ctx, nodeID)
}
