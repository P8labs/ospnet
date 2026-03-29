package handlers

import (
	"encoding/json"
	"net/http"

	"ospnet/internal/agent/config"
	"ospnet/internal/agent/runtime/manager"
	"ospnet/pkg/res"
)

type Handlers struct {
	nodeCfg  config.NodeConfig
	manager  *manager.Manager
}

type RunRequest struct {
	Image string `json:"image"`
	Name  string `json:"name"`
	Port  int    `json:"port"`
}

type StopRequest struct {
	Name string `json:"name"`
}

func New(nodeCfg config.NodeConfig, mgr *manager.Manager) *Handlers {
	return &Handlers{nodeCfg: nodeCfg, manager: mgr}
}

func (h *Handlers) Register(mux *http.ServeMux) {
	mux.HandleFunc("/health", h.handleHealth)
	mux.HandleFunc("/containers", h.handleListContainers)
	mux.HandleFunc("/containers/run", h.handleRunContainer)
	mux.HandleFunc("/containers/stop", h.handleStopContainer)
}

func (h *Handlers) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		res.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	res.SuccessRaw(w, res.M{
		"status":  "ok",
		"node_id": h.nodeCfg.NodeID,
		"ip":      h.nodeCfg.IP,
	})
}

func (h *Handlers) handleRunContainer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		res.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RunRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		res.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	record, err := h.manager.StartContainer(r.Context(), manager.StartRequest{
		Image: req.Image,
		Name:  req.Name,
		Port:  req.Port,
	})
	if err != nil {
		res.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Success(w, res.M{"status": "started", "id": record.ID, "docker_id": record.DockerID})
}

func (h *Handlers) handleStopContainer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		res.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req StopRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		res.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if err := h.manager.StopContainer(r.Context(), req.Name); err != nil {
		res.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Success(w, res.M{"status": "stopped", "name": req.Name})
}

func (h *Handlers) handleListContainers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		res.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	containers, err := h.manager.GetContainers(r.Context())
	if err != nil {
		res.Error(w, "failed to list containers", http.StatusInternalServerError)
		return
	}

	res.Success(w, containers)
}
