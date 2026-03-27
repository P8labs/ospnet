package master

import (
	"encoding/json"
	"net/http"
	"ospnet/pkg/types"
)

type RegisterRequest struct {
	NodeID string `json:"node_id"`
	IP     string `json:"ip"`
	CPU    int    `json:"cpu"`
	Memory int    `json:"memory"`
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	node := &types.Node{
		ID:     req.NodeID,
		IP:     req.IP,
		CPU:    req.CPU,
		Memory: req.Memory,
	}

	s.store.Register(node)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"registered"}`))
}

type HeartbeatRequest struct {
	NodeID string `json:"node_id"`
}

func (s *Server) handleHeartbeat(w http.ResponseWriter, r *http.Request) {
	var req HeartbeatRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	s.store.Heartbeat(req.NodeID)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func (s *Server) handleListNodes(w http.ResponseWriter, r *http.Request) {
	nodes := s.store.List()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nodes)
}
