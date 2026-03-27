package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"ospnet/internal/master/config"
	"ospnet/pkg/res"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
)

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	c := cors.AllowAll()
	r.Use(c.Handler)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		if err := app.db.Ping(r.Context()); err != nil {
			res.Error(w, "Database unavailable", 503)
			return
		}

		res.Success(w, map[string]string{
			"status":         "ok",
			"timestamp":      time.Now().Format(time.RFC3339),
			"server_version": version,
		}, "OK")
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, app.config.FrontendURL, http.StatusTemporaryRedirect)
	})

	return r
}

func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", app.config.Port),
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	app.logger.Info("server has started at addr", "port", app.config.Port)
	return srv.ListenAndServe()
}

type application struct {
	config *config.Config
	logger *slog.Logger
	db     *pgxpool.Pool
}
