package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	DefaultConfigPath        = "/etc/ospnet/config.json"
	DefaultTokenPath         = "/etc/ospnet/token"
	DefaultDBPath            = "/var/lib/ospnet/agent.db"
	DefaultPort              = 9000
	DefaultHeartbeatInterval = 12 * time.Second
)

type RuntimeConfig struct {
	ConfigPath        string
	TokenPath         string
	DBPath            string
	MasterURL         string
	Port              int
	BindAddr          string
	NodeName          string
	NodeRegion        string
	NodeType          string
	HeartbeatInterval time.Duration
}

type NodeConfig struct {
	NodeID    string            `json:"node_id"`
	NodeName  string            `json:"node_name"`
	MasterURL string            `json:"master_url"`
	Labels    []string          `json:"labels"`
	Metadata  map[string]string `json:"metadata"`
	Hostname  string            `json:"hostname"`
	CPU       int64             `json:"cpu"`
	Memory    int64             `json:"memory"`
	Arch      string            `json:"arch"`
	IP        string            `json:"ip"`
}

type FileStore struct {
	path string
}

func NewFileStore(path string) *FileStore {
	return &FileStore{path: path}
}

func (f *FileStore) Load() (NodeConfig, error) {
	raw, err := os.ReadFile(f.path)
	if err != nil {
		return NodeConfig{}, err
	}

	if len(raw) == 0 {
		return NodeConfig{}, nil
	}

	var cfg NodeConfig
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return NodeConfig{}, fmt.Errorf("failed to parse node config: %w", err)
	}

	return cfg, nil
}

func (f *FileStore) Save(cfg NodeConfig) error {
	if err := os.MkdirAll(filepath.Dir(f.path), 0o755); err != nil {
		return fmt.Errorf("failed to create config dir: %w", err)
	}

	raw, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}

	tmpPath := f.path + ".tmp"
	if err := os.WriteFile(tmpPath, raw, 0o600); err != nil {
		return fmt.Errorf("failed to write temp config: %w", err)
	}

	if err := os.Rename(tmpPath, f.path); err != nil {
		return fmt.Errorf("failed to replace config: %w", err)
	}

	return nil
}

func LoadRuntimeConfigFromEnv() RuntimeConfig {
	return RuntimeConfig{
		ConfigPath:        getenvOrDefault("OSPNET_CONFIG_PATH", DefaultConfigPath),
		TokenPath:         getenvOrDefault("OSPNET_TOKEN_PATH", DefaultTokenPath),
		DBPath:            getenvOrDefault("OSPNET_DB_PATH", DefaultDBPath),
		MasterURL:         "http://localhost:8000",
		Port:              parseIntOrDefault(os.Getenv("OSPNET_AGENT_PORT"), DefaultPort),
		BindAddr:          strings.TrimSpace(os.Getenv("OSPNET_BIND_ADDR")),
		NodeName:          strings.TrimSpace(os.Getenv("OSPNET_NODE_NAME")),
		NodeRegion:        strings.TrimSpace(os.Getenv("OSPNET_NODE_REGION")),
		NodeType:          strings.TrimSpace(os.Getenv("OSPNET_NODE_TYPE")),
		HeartbeatInterval: parseDurationOrDefault(os.Getenv("OSPNET_HEARTBEAT_INTERVAL"), DefaultHeartbeatInterval),
	}
}

func getenvOrDefault(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func parseIntOrDefault(raw string, fallback int) int {
	value := strings.TrimSpace(raw)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func parseDurationOrDefault(raw string, fallback time.Duration) time.Duration {
	value := strings.TrimSpace(raw)
	if value == "" {
		return fallback
	}
	parsed, err := time.ParseDuration(value)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}

func EnsurePaths(cfg RuntimeConfig) error {
	paths := []string{
		cfg.ConfigPath,
		cfg.TokenPath,
		cfg.DBPath,
	}

	for _, p := range paths {
		dir := filepath.Dir(p)

		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("failed to create dir for %s: %w", p, err)
		}

		if _, err := os.Stat(p); os.IsNotExist(err) {
			file, err := os.OpenFile(p, os.O_CREATE, 0o600)
			if err != nil {
				return fmt.Errorf("failed to create file %s: %w", p, err)
			}
			file.Close()
		}
	}

	return nil
}
