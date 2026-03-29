package network

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"
)

type Tailscale struct {
	logger *log.Logger
}

func NewTailscale(logger *log.Logger) *Tailscale {
	return &Tailscale{logger: logger}
}

func (t *Tailscale) UpWithAuthKey(ctx context.Context, authKey string) error {
	if strings.TrimSpace(authKey) == "" {
		return fmt.Errorf("tailscale auth key cannot be empty")
	}

	var lastErr error
	for attempt := 1; attempt <= 5; attempt++ {
		cmdCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		cmd := exec.CommandContext(cmdCtx, "tailscale", "up", "--authkey="+authKey)
		output, err := cmd.CombinedOutput()
		cancel()

		if err == nil {
			return nil
		}

		lastErr = fmt.Errorf("tailscale up attempt %d failed: %s", attempt, strings.TrimSpace(string(output)))
		t.logger.Printf("%v", lastErr)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Duration(attempt*2) * time.Second):
		}
	}

	return lastErr
}

func (t *Tailscale) IPv4(ctx context.Context) (string, error) {
	cmdCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(cmdCtx, "tailscale", "ip", "-4")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get tailscale ipv4: %s", strings.TrimSpace(string(output)))
	}

	for _, line := range strings.Split(string(output), "\n") {
		value := strings.TrimSpace(line)
		if value != "" {
			return value, nil
		}
	}

	return "", fmt.Errorf("tailscale returned no ipv4")
}
