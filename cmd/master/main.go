package main

import (
	"context"
	"log/slog"
	"os"

	"ospnet/internal/master/config"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var version = "dev"

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	err := godotenv.Load()
	if err != nil {
		logger.Info(
			"No .env file loaded, falling back to system environment",
			slog.String("error", err.Error()),
		)
	}

	ctx := context.Background()

	cfg := config.Load()

	logger.Info("[OSPNet MASTER]: starting server", "version", version)

	dbpool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer dbpool.Close()

	if err := dbpool.Ping(ctx); err != nil {
		panic(err)
	}

	api := application{
		config: cfg,
		db:     dbpool,
		logger: logger,
	}
	if err := api.run(api.mount()); err != nil {
		logger.Error("[OSPNet MASTER]: server failed to start", "error", err)
		os.Exit(1)
	}
}
