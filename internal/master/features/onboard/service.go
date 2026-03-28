package onboard

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"ospnet/internal/master/config"
	"ospnet/internal/master/core"
	"ospnet/internal/master/db"
	"ospnet/internal/master/tailscale"

	"github.com/google/uuid"
)

type OnboardService interface {
	CreateToken(ctx context.Context) (CreateTokenResponse, error)
	RegisterNode(ctx context.Context, req RegisterNodeRequest) (RegisterNodeResponse, error)
}

type sv struct {
	db  db.Querier
	cfg config.Config
	ts  tailscale.TailscaleClient
}

func NewService(db db.Querier, cfg config.Config, ts tailscale.TailscaleClient) OnboardService {
	return &sv{db: db, cfg: cfg, ts: ts}
}

func (s *sv) CreateToken(ctx context.Context) (CreateTokenResponse, error) {
	token := uuid.New().String()
	expiresAt := time.Now().Add(10 * time.Minute)

	err := s.db.CreateToken(ctx, db.CreateTokenParams{
		Token:     token,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return CreateTokenResponse{}, err
	}

	return CreateTokenResponse{
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}

func (s *sv) RegisterNode(ctx context.Context, req RegisterNodeRequest) (RegisterNodeResponse, error) {
	if req.Token == "" || req.NodeID == "" || req.IP == "" {
		return RegisterNodeResponse{}, errors.New("missing required fields")
	}

	req.Type = strings.ToLower(strings.TrimSpace(req.Type))
	req.Region = strings.ToLower(strings.TrimSpace(req.Region))

	tokenRow, err := s.db.GetValidToken(ctx, req.Token)
	if err != nil {
		return RegisterNodeResponse{}, fmt.Errorf("invalid or expired token")
	}

	labels := []string{
		"tag:ospnet-node",
	}

	if req.Type != "" {
		labels = append(labels, "tag:"+req.Type)
	}
	if req.Region != "" {
		labels = append(labels, "tag:"+req.Region)
	}

	authKey, err := s.ts.CreateAuthKey(tailscale.CreateKeyOptions{
		Tags:      labels,
		Ephemeral: false,
		Reusable:  false,
		Expiry:    1800,
	})
	if err != nil {
		return RegisterNodeResponse{}, fmt.Errorf("failed to create auth key: %w", err)
	}

	err = s.db.CreateNode(ctx, db.CreateNodeParams{
		ID:       req.NodeID,
		Name:     req.Name,
		Hostname: req.Hostname,
		Ip:       req.IP,
		Cpu:      req.CPU,
		Memory:   req.Memory,
		Arch:     req.Arch,
		Region:   req.Region,
		Type:     req.Type,
		Status:   "healthy",
		LastSeen: time.Now(),
	})
	if err != nil {
		return RegisterNodeResponse{}, err
	}

	err = s.db.MarkTokenUsed(ctx, tokenRow.Token)
	if err != nil {
		return RegisterNodeResponse{}, err
	}

	return RegisterNodeResponse{
		AuthKey:  authKey,
		NodeName: "ospnet-" + core.SanitizeName(req.Name),
		Labels:   labels,
	}, nil
}
