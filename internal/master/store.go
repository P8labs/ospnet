package master

import (
	"ospnet/pkg/types"
	"sync"
	"time"
)

type Store struct {
	mu    sync.RWMutex
	nodes map[string]*types.Node
}

func NewStore() *Store {
	return &Store{
		nodes: make(map[string]*types.Node),
	}
}

func (s *Store) Register(node *types.Node) {
	s.mu.Lock()
	defer s.mu.Unlock()

	node.LastSeen = time.Now()
	node.Status = "healthy"

	s.nodes[node.ID] = node
}

func (s *Store) Heartbeat(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if node, ok := s.nodes[id]; ok {
		node.LastSeen = time.Now()
		node.Status = "healthy"
	}
}

func (s *Store) List() []*types.Node {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]*types.Node, 0, len(s.nodes))
	for _, n := range s.nodes {
		out = append(out, n)
	}
	return out
}
