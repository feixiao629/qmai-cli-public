package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigFrom_NonExistentFile(t *testing.T) {
	cfg, err := LoadConfigFrom("/tmp/qmai-test-nonexistent/config.yaml")
	if err != nil {
		t.Fatalf("expected no error for non-existent file, got: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
	if len(cfg.Profiles) != 0 {
		t.Errorf("expected empty profiles, got %d", len(cfg.Profiles))
	}
	if cfg.ActiveProfile != "" {
		t.Errorf("expected empty active profile, got %q", cfg.ActiveProfile)
	}
}

func TestSaveAndLoadRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	cfg := &CliConfig{
		ActiveProfile: "prod",
		DefaultFormat: "json",
		Debug:         true,
		Profiles:      make(map[string]*Profile),
		path:          path,
	}
	cfg.SetProfile("prod", &Profile{
		ShopCode: "store-123",
		BaseURL: "https://example.com/api/",
	})

	if err := cfg.Save(); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Verify the file was created
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("config file not created: %v", err)
	}

	loaded, err := LoadConfigFrom(path)
	if err != nil {
		t.Fatalf("LoadConfigFrom failed: %v", err)
	}

	if loaded.ActiveProfile != "prod" {
		t.Errorf("ActiveProfile = %q, want %q", loaded.ActiveProfile, "prod")
	}
	if loaded.DefaultFormat != "json" {
		t.Errorf("DefaultFormat = %q, want %q", loaded.DefaultFormat, "json")
	}
	if !loaded.Debug {
		t.Error("Debug = false, want true")
	}

	p := loaded.Profiles["prod"]
	if p == nil {
		t.Fatal("profile 'prod' not found")
	}
	if p.Name != "prod" {
		t.Errorf("profile Name = %q, want %q", p.Name, "prod")
	}
	if p.ShopCode != "store-123" {
		t.Errorf("profile ShopCode = %q, want %q", p.ShopCode, "store-123")
	}
	if p.BaseURL != "https://example.com/api/" {
		t.Errorf("profile BaseURL = %q, want %q", p.BaseURL, "https://example.com/api/")
	}
}

func TestActiveProfileConfig(t *testing.T) {
	cfg := &CliConfig{
		Profiles: make(map[string]*Profile),
	}

	// No active profile set
	if p := cfg.ActiveProfileConfig(); p != nil {
		t.Errorf("expected nil when no active profile, got %+v", p)
	}

	// Active profile set but doesn't exist in map
	cfg.ActiveProfile = "missing"
	if p := cfg.ActiveProfileConfig(); p != nil {
		t.Errorf("expected nil for missing profile, got %+v", p)
	}

	// Active profile exists
	cfg.SetProfile("dev", &Profile{ShopCode: "s1"})
	cfg.ActiveProfile = "dev"
	p := cfg.ActiveProfileConfig()
	if p == nil {
		t.Fatal("expected non-nil profile")
	}
	if p.ShopCode != "s1" {
		t.Errorf("ShopCode = %q, want %q", p.ShopCode, "s1")
	}
}

func TestBaseURL_Defaults(t *testing.T) {
	cfg := &CliConfig{
		Profiles: make(map[string]*Profile),
	}

	// No active profile -> default
	if got := cfg.BaseURL(); got != DefaultBaseURL {
		t.Errorf("BaseURL() = %q, want %q", got, DefaultBaseURL)
	}

	// Active profile with no BaseURL -> default
	cfg.SetProfile("test", &Profile{})
	cfg.ActiveProfile = "test"
	if got := cfg.BaseURL(); got != DefaultBaseURL {
		t.Errorf("BaseURL() = %q, want %q", got, DefaultBaseURL)
	}

	// Active profile with custom BaseURL
	cfg.Profiles["test"].BaseURL = "https://custom.api/"
	if got := cfg.BaseURL(); got != "https://custom.api/" {
		t.Errorf("BaseURL() = %q, want %q", got, "https://custom.api/")
	}
}

func TestFormat_Defaults(t *testing.T) {
	cfg := &CliConfig{}

	// No format set -> "table"
	if got := cfg.Format(); got != "table" {
		t.Errorf("Format() = %q, want %q", got, "table")
	}

	// Custom format
	cfg.DefaultFormat = "csv"
	if got := cfg.Format(); got != "csv" {
		t.Errorf("Format() = %q, want %q", got, "csv")
	}
}

func TestSetProfile(t *testing.T) {
	cfg := &CliConfig{
		Profiles: make(map[string]*Profile),
	}

	cfg.SetProfile("staging", &Profile{
		ShopCode: "store-staging",
		BaseURL: "https://staging.api/",
	})

	p, ok := cfg.Profiles["staging"]
	if !ok {
		t.Fatal("profile 'staging' not found after SetProfile")
	}
	if p.Name != "staging" {
		t.Errorf("Name = %q, want %q", p.Name, "staging")
	}
	if p.ShopCode != "store-staging" {
		t.Errorf("ShopCode = %q, want %q", p.ShopCode, "store-staging")
	}

	// Overwrite existing profile
	cfg.SetProfile("staging", &Profile{
		ShopCode: "store-staging-v2",
	})
	p = cfg.Profiles["staging"]
	if p.ShopCode != "store-staging-v2" {
		t.Errorf("ShopCode after overwrite = %q, want %q", p.ShopCode, "store-staging-v2")
	}
}
