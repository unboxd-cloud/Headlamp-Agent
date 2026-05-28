package nodeagent

import (
	"encoding/json"
	"os"
)

type Config struct {
	AgentID string `json:"agentId"`
	HostID  string `json:"hostId"`
	Mode    string `json:"mode"`
	Listen  string `json:"listen"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if cfg.Listen == "" {
		cfg.Listen = "127.0.0.1:9080"
	}

	return &cfg, nil
}
