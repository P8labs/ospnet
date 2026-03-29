package onboard

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"

	"ospnet/internal/agent/config"
	"ospnet/internal/agent/system"
)

type ConfigStore interface {
	Load() (config.NodeConfig, error)
	Save(config.NodeConfig) error
}

type SystemCollector interface {
	Collect() (system.Info, error)
}

type TailscaleClient interface {
	UpWithAuthKey(ctx context.Context, authKey string) error
	IPv4(ctx context.Context) (string, error)
}

type Service struct {
	runtimeCfg config.RuntimeConfig
	cfgStore   ConfigStore
	system     SystemCollector
	tailscale  TailscaleClient
	httpClient *http.Client
	logger     *log.Logger
}

func NewService(runtimeCfg config.RuntimeConfig, cfgStore ConfigStore, sys SystemCollector, ts TailscaleClient, logger *log.Logger) *Service {
	return &Service{
		runtimeCfg: runtimeCfg,
		cfgStore:   cfgStore,
		system:     sys,
		tailscale:  ts,
		httpClient: &http.Client{},
		logger:     logger,
	}
}

type registerNodeRequest struct {
	Token    string `json:"token"`
	NodeID   string `json:"node_id"`
	Hostname string `json:"hostname"`
	Name     string `json:"name"`
	Region   string `json:"region"`
	Type     string `json:"type"`
	CPU      int64  `json:"cpu"`
	Memory   int64  `json:"memory"`
	Arch     string `json:"arch"`
	IP       string `json:"ip"`
}

type registerNodeResponse struct {
	AuthKey  string   `json:"auth_key"`
	NodeName string   `json:"node_name"`
	Labels   []string `json:"labels"`
}

func (s *Service) EnsureOnboarded(ctx context.Context) (config.NodeConfig, error) {
	nodeCfg, err := s.cfgStore.Load()
	if err == nil {
		if nodeCfg.MasterURL == "" {
			nodeCfg.MasterURL = s.runtimeCfg.MasterURL
		}
		s.logger.Printf("existing node config found, onboarding skipped")
		return nodeCfg, nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return config.NodeConfig{}, fmt.Errorf("failed to load node config: %w", err)
	}

	tokenRaw, err := os.ReadFile(s.runtimeCfg.TokenPath)
	if err != nil {
		return config.NodeConfig{}, fmt.Errorf("failed to read onboarding token: %w", err)
	}
	token := strings.TrimSpace(string(tokenRaw))
	if token == "" {
		return config.NodeConfig{}, fmt.Errorf("onboarding token is empty")
	}

	systemInfo, err := s.system.Collect()
	if err != nil {
		return config.NodeConfig{}, fmt.Errorf("failed to collect system info: %w", err)
	}

	name := s.runtimeCfg.NodeName
	if name == "" {
		name = systemInfo.Hostname
	}

	nodeID := uuid.NewString()
	registerReq := registerNodeRequest{
		Token:    token,
		NodeID:   nodeID,
		Hostname: systemInfo.Hostname,
		Name:     name,
		Region:   s.runtimeCfg.NodeRegion,
		Type:     s.runtimeCfg.NodeType,
		CPU:      systemInfo.CPU,
		Memory:   systemInfo.Memory,
		Arch:     systemInfo.Arch,
		IP:       system.FirstPrivateIPv4(),
	}

	registerResp, err := s.registerWithMaster(ctx, registerReq)
	if err != nil {
		return config.NodeConfig{}, err
	}

	if err := s.tailscale.UpWithAuthKey(ctx, registerResp.AuthKey); err != nil {
		return config.NodeConfig{}, fmt.Errorf("tailscale join failed: %w", err)
	}

	tsIP, err := s.tailscale.IPv4(ctx)
	if err != nil {
		s.logger.Printf("failed to get tailscale ipv4, continuing without it: %v", err)
	}

	nodeCfg = config.NodeConfig{
		NodeID:    nodeID,
		NodeName:  registerResp.NodeName,
		MasterURL: s.runtimeCfg.MasterURL,
		Labels:    registerResp.Labels,
		Metadata: map[string]string{
			"name":   name,
			"region": registerReq.Region,
			"type":   registerReq.Type,
		},
		Hostname: systemInfo.Hostname,
		CPU:      systemInfo.CPU,
		Memory:   systemInfo.Memory,
		Arch:     systemInfo.Arch,
		IP:       tsIP,
	}
	if nodeCfg.NodeName == "" {
		nodeCfg.NodeName = systemInfo.Hostname
	}

	if err := s.cfgStore.Save(nodeCfg); err != nil {
		return config.NodeConfig{}, fmt.Errorf("failed to save node config: %w", err)
	}

	if err := os.Remove(s.runtimeCfg.TokenPath); err != nil && !errors.Is(err, os.ErrNotExist) {
		s.logger.Printf("warning: failed to remove token file: %v", err)
	}

	return nodeCfg, nil
}

func (s *Service) registerWithMaster(ctx context.Context, payload registerNodeRequest) (registerNodeResponse, error) {
	raw, err := json.Marshal(payload)
	if err != nil {
		return registerNodeResponse{}, fmt.Errorf("failed to encode register request: %w", err)
	}

	path := "/api/onboard/register"
	url := s.runtimeCfg.MasterURL + path
	responseBody, err := s.postJSON(ctx, url, raw)
	if err != nil {
		return registerNodeResponse{}, err
	}

	response, err := parseRegisterResponse(responseBody)
	if err != nil {
		return registerNodeResponse{}, err

	}
	if strings.TrimSpace(response.AuthKey) == "" {
		return registerNodeResponse{}, err

	}
	if strings.TrimSpace(response.NodeName) == "" {
		response.NodeName = payload.Hostname
	}

	return response, nil

}

func (s *Service) postJSON(ctx context.Context, endpoint string, payload []byte) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("master returned %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	return body, nil
}

func parseRegisterResponse(raw []byte) (registerNodeResponse, error) {
	var direct registerNodeResponse
	if err := json.Unmarshal(raw, &direct); err == nil && direct.AuthKey != "" {
		return direct, nil
	}

	type wrappedResponse struct {
		Data registerNodeResponse `json:"data"`
	}
	var wrapped wrappedResponse
	if err := json.Unmarshal(raw, &wrapped); err != nil {
		return registerNodeResponse{}, err
	}
	if wrapped.Data.AuthKey == "" {
		return registerNodeResponse{}, fmt.Errorf("invalid response structure")
	}
	return wrapped.Data, nil
}
