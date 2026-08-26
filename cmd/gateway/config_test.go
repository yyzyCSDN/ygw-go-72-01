package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfigValidate(t *testing.T) {
	cfg := DefaultConfig()
	if err := cfg.Validate(); err != nil {
		t.Fatalf("default config invalid: %v", err)
	}
}

func TestLoadConfigFromFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "gateway.json")
	content := `{"addr":"127.0.0.1:9000","cycle_seconds":7,"rotate_seconds":90}`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if cfg.Addr != "127.0.0.1:9000" {
		t.Fatalf("addr = %q", cfg.Addr)
	}
	if cfg.CycleSeconds != 7 {
		t.Fatalf("cycle = %d", cfg.CycleSeconds)
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("loaded config invalid: %v", err)
	}
}

func TestLoadConfigEmptyPath(t *testing.T) {
	cfg, err := LoadConfig("")
	if err != nil {
		t.Fatalf("load empty failed: %v", err)
	}
	if cfg.Addr != DefaultConfig().Addr {
		t.Fatalf("default addr = %q", cfg.Addr)
	}
}

func TestLoadConfigMissingFile(t *testing.T) {
	if _, err := LoadConfig(filepath.Join(t.TempDir(), "missing.json")); err == nil {
		t.Fatal("expected error for missing config")
	}
}

func TestConfigMerge(t *testing.T) {
	base := DefaultConfig()
	merged := base.Merge(Config{Addr: "0.0.0.0:8123", CycleSeconds: 9})
	if merged.Addr != "0.0.0.0:8123" {
		t.Fatalf("merged addr = %q", merged.Addr)
	}
	if merged.CycleSeconds != 9 {
		t.Fatalf("merged cycle = %d", merged.CycleSeconds)
	}
	if merged.RotateSeconds != base.RotateSeconds {
		t.Fatalf("rotate changed unexpectedly: %d", merged.RotateSeconds)
	}
}

func TestConfigValidateRejectsZeroCycle(t *testing.T) {
	cfg := DefaultConfig()
	cfg.CycleSeconds = 0
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected invalid config")
	}
}
