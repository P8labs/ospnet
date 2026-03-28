package onboard

import (
	"ospnet/internal/master/config"
	"ospnet/internal/master/db"
	"ospnet/internal/master/tailscale"

	"github.com/go-chi/chi/v5"
)

func Routes(db db.Querier, cfg config.Config) *chi.Mux {
	r := chi.NewRouter()

	ts := tailscale.NewClient(cfg.TsAuthKey, cfg.TailnetID)

	service := NewService(db, cfg, ts)
	handler := NewHandler(service)

	r.Post("/onboard/init", handler.CreateToken)
	r.Post("/onboard/register", handler.RegisterNode)
	return r
}
