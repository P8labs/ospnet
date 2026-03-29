package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"ospnet/internal/agent/app"
	"ospnet/internal/agent/config"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("No .env file loaded (%v), falling back to system environment", err)
	}
	logger := log.New(os.Stdout, "[ospnet-agent] ", log.LstdFlags|log.LUTC)

	rootCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	runtimeCfg := config.LoadRuntimeConfigFromEnv()
	if err := config.EnsurePaths(runtimeCfg); err != nil {
		panic(err)
	}
	agentApp := app.New(runtimeCfg, logger)

	if err := agentApp.Run(rootCtx); err != nil {
		logger.Fatalf("agent stopped with error: %v", err)
	}
}
