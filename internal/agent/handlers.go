package agent

import (
	"encoding/json"
	"net/http"
	"ospnet/pkg/res"
)

type RunRequest struct {
	Image string `json:"image"`
	Name  string `json:"name"`
	Port  int    `json:"port"`
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	res.SuccessRaw(w, res.M{"status": "ok"})
}

func (s *Server) handleRunContainer(w http.ResponseWriter, r *http.Request) {
	var req RunRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		res.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if err := RunContainer(req); err != nil {
		res.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Success(w, res.M{"status": "started"})
}
