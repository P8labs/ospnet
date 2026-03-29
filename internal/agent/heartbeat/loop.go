package heartbeat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"ospnet/internal/agent/config"
)

type Loop struct {
	node     config.NodeConfig
	interval time.Duration
	client   *http.Client
	logger   *log.Logger
}

func New(node config.NodeConfig, interval time.Duration, logger *log.Logger) *Loop {
	if interval <= 0 {
		interval = config.DefaultHeartbeatInterval
	}
	return &Loop{
		node:     node,
		interval: interval,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		logger: logger,
	}
}

func (l *Loop) Run(ctx context.Context) error {
	ticker := time.NewTicker(l.interval)
	defer ticker.Stop()

	l.send(ctx)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			l.send(ctx)
		}
	}
}

func (l *Loop) send(ctx context.Context) {
	payload := map[string]string{
		"node_id": l.node.NodeID,
		"ip":      l.node.IP,
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		l.logger.Printf("heartbeat marshal failed: %v", err)
		return
	}

	paths := []string{"/nodes/heartbeat", "/api/nodes/heartbeat"}
	var lastErr error
	for _, path := range paths {
		url := l.node.MasterURL + path
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(raw))
		if err != nil {
			lastErr = err
			continue
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := l.client.Do(req)
		if err != nil {
			lastErr = err
			continue
		}
		body, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return
		}
		lastErr = fmt.Errorf("status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	if lastErr != nil {
		l.logger.Printf("heartbeat failed: %v", lastErr)
	}
}
