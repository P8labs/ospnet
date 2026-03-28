package node

import (
	"encoding/json"
	"net/http"
	"ospnet/pkg/res"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) ListNodes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	nodes, err := h.svc.ListNodes(ctx)
	if err != nil {
		res.Error(w, "failed to fetch nodes", http.StatusInternalServerError)
		return
	}

	res.Success(w, nodes)
}

type HeartbeatRequest struct {
	NodeID string `json:"node_id"`
}

func (h *Handler) Heartbeat(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req HeartbeatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		res.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.NodeID == "" {
		res.Error(w, "node_id required", http.StatusBadRequest)
		return
	}

	if err := h.svc.Heartbeat(ctx, req.NodeID); err != nil {
		res.Error(w, "failed to update heartbeat", http.StatusInternalServerError)
		return
	}

	res.Success(w, map[string]string{
		"status": "ok",
	})
}
