package onboard

import (
	"net/http"
	"ospnet/pkg/res"
)

type OnboardHandler struct {
	svc OnboardService
}

func NewHandler(svc OnboardService) *OnboardHandler {
	return &OnboardHandler{svc: svc}
}

func (h *OnboardHandler) CreateToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	resp, err := h.svc.CreateToken(ctx)
	if err != nil {
		res.Error(w, "failed to create token", http.StatusInternalServerError)
		return
	}

	res.Success(w, resp)
}

func (h *OnboardHandler) RegisterNode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req, err := res.Parse[RegisterNodeRequest](r)
	if err != nil {
		res.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.svc.RegisterNode(ctx, req)
	if err != nil {
		res.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res.Success(w, resp)
}
