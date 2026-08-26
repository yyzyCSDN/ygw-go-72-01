package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Addr          string `json:"addr"`
	DataDir       string `json:"data_dir"`
	WebDir        string `json:"web_dir"`
	CycleSeconds  int    `json:"cycle_seconds"`
	RotateSeconds int    `json:"rotate_seconds"`
}

func DefaultConfig() Config {
	return Config{
		Addr:          "127.0.0.1:8090",
		DataDir:       "data",
		WebDir:        "web",
		CycleSeconds:  3,
		RotateSeconds: 60,
	}
}

func LoadConfig(path string) (Config, error) {
	if path == "" {
		return DefaultConfig(), nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (c Config) Validate() error {
	if c.Addr == "" {
		return fmt.Errorf("addr must not be empty")
	}
	if c.CycleSeconds <= 0 {
		return fmt.Errorf("cycle_seconds must be positive")
	}
	if c.RotateSeconds <= 0 {
		return fmt.Errorf("rotate_seconds must be positive")
	}
	return nil
}

func (c Config) Merge(overlay Config) Config {
	if overlay.Addr != "" {
		c.Addr = overlay.Addr
	}
	if overlay.DataDir != "" {
		c.DataDir = overlay.DataDir
	}
	if overlay.WebDir != "" {
		c.WebDir = overlay.WebDir
	}
	if overlay.CycleSeconds > 0 {
		c.CycleSeconds = overlay.CycleSeconds
	}
	if overlay.RotateSeconds > 0 {
		c.RotateSeconds = overlay.RotateSeconds
	}
	return c
}
