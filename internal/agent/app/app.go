package app

import (
	"context"
	"errors"
	"fmt"
	"log"

	"ospnet/internal/agent/config"
	types "ospnet/internal/agent/db"
	"ospnet/internal/agent/handlers"
	"ospnet/internal/agent/heartbeat"
	"ospnet/internal/agent/network"
	"ospnet/internal/agent/onboard"
	"ospnet/internal/agent/reconciler"
	runtimedocker "ospnet/internal/agent/runtime/docker"
	"ospnet/internal/agent/runtime/manager"
	"ospnet/internal/agent/server"
	"ospnet/internal/agent/system"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type App struct {
	runtimeCfg config.RuntimeConfig
	logger     *log.Logger
}

func New(runtimeCfg config.RuntimeConfig, logger *log.Logger) *App {
	return &App{runtimeCfg: runtimeCfg, logger: logger}
}

func (a *App) Run(ctx context.Context) error {

	configStore := config.NewFileStore(a.runtimeCfg.ConfigPath)
	db, err := gorm.Open(sqlite.Open(a.runtimeCfg.DBPath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&types.Containers{})

	onboardSvc := onboard.NewService(
		a.runtimeCfg,
		configStore,
		system.NewCollector(),
		network.NewTailscale(a.logger),
		a.logger,
	)

	nodeCfg, err := onboardSvc.EnsureOnboarded(ctx)
	if err != nil {
		return fmt.Errorf("onboarding failed: %w", err)
	}

	dockerClient, err := runtimedocker.New()
	if err != nil {
		return fmt.Errorf("failed to initialize docker client: %w", err)
	}

	containerManager := manager.New(db, dockerClient, a.logger)
	handler := handlers.New(nodeCfg, containerManager)

	bindAddr := a.runtimeCfg.BindAddr
	if bindAddr == "" {
		if nodeCfg.IP != "" {
			bindAddr = nodeCfg.IP
		} else {
			bindAddr = "127.0.0.1"
		}
	}

	httpServer := server.New(bindAddr, a.runtimeCfg.Port, handler, a.logger)
	heartbeatLoop := heartbeat.New(nodeCfg, a.runtimeCfg.HeartbeatInterval, a.logger)
	reconcilerLoop := reconciler.New(containerManager, a.logger)

	errCh := make(chan error, 3)
	go func() { errCh <- httpServer.Run(ctx) }()
	go func() { errCh <- heartbeatLoop.Run(ctx) }()
	go func() { errCh <- reconcilerLoop.Run(ctx) }()

	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-errCh:
			if err != nil && !errors.Is(err, context.Canceled) {
				return err
			}
		}
	}
}
