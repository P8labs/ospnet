package node

import (
	"ospnet/internal/master/config"
	"ospnet/internal/master/db"

	"github.com/go-chi/chi/v5"
)

func Routes(db db.Querier, cfg config.Config) *chi.Mux {
	r := chi.NewRouter()

	service := NewService(db)
	handler := NewHandler(service)

	r.Get("/", handler.ListNodes)
	r.Post("/heartbeat", handler.Heartbeat)
	return r
}
