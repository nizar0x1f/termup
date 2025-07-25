package config

import (
	"os"
	"testing"
)

func TestSimpleConfig(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "testconfig")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	cfg := &Config{
		AccessKeyID:     "test-access-key",
		SecretAccessKey: "test-secret-key",
		Bucket:          "test-bucket",
		Endpoint:        "test-endpoint",
	}

	if err := Save(cfg); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	loadedCfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if loadedCfg.AccessKeyID != cfg.AccessKeyID {
		t.Errorf("AccessKeyID mismatch: got %v, want %v", loadedCfg.AccessKeyID, cfg.AccessKeyID)
	}
}
