package cmdutil

import (
	"testing"

	"github.com/madaima/qmai-cli/internal/core"
)

func TestNewFactory_ReturnsNonNil(t *testing.T) {
	f := NewFactory()
	if f == nil {
		t.Fatal("NewFactory returned nil")
	}
	if f.IOStreams == nil {
		t.Error("IOStreams is nil")
	}
	if f.ConfigFunc == nil {
		t.Error("ConfigFunc is nil")
	}
}

func TestEffectiveFormat_WithFlag(t *testing.T) {
	f := NewFactory()
	// Override ConfigFunc so it doesn't try to load from disk
	f.ConfigFunc = func() (*core.CliConfig, error) {
		return &core.CliConfig{DefaultFormat: "csv"}, nil
	}

	// When the flag is set, it should take precedence over config
	f.Format = "json"
	if got := f.EffectiveFormat(); got != "json" {
		t.Errorf("EffectiveFormat() = %q, want %q", got, "json")
	}
}

func TestEffectiveFormat_WithoutFlag_UsesConfig(t *testing.T) {
	f := &Factory{
		IOStreams: DefaultIOStreams(),
	}
	f.ConfigFunc = func() (*core.CliConfig, error) {
		return &core.CliConfig{DefaultFormat: "csv"}, nil
	}

	// No flag set -> should use config format
	if got := f.EffectiveFormat(); got != "csv" {
		t.Errorf("EffectiveFormat() = %q, want %q", got, "csv")
	}
}

func TestEffectiveFormat_NoFlagNoConfig_DefaultsToTable(t *testing.T) {
	f := &Factory{
		IOStreams: DefaultIOStreams(),
	}
	f.ConfigFunc = func() (*core.CliConfig, error) {
		// Config with no DefaultFormat set -> Format() returns "table"
		return &core.CliConfig{}, nil
	}

	if got := f.EffectiveFormat(); got != "table" {
		t.Errorf("EffectiveFormat() = %q, want %q", got, "table")
	}
}

func TestConfig_LazyLoading(t *testing.T) {
	callCount := 0
	expectedCfg := &core.CliConfig{ActiveProfile: "test"}

	f := &Factory{
		IOStreams: DefaultIOStreams(),
	}
	f.ConfigFunc = func() (*core.CliConfig, error) {
		callCount++
		return expectedCfg, nil
	}

	// First call
	cfg1, err := f.Config()
	if err != nil {
		t.Fatalf("Config() error: %v", err)
	}
	if cfg1 != expectedCfg {
		t.Error("Config() returned unexpected config")
	}
	if callCount != 1 {
		t.Errorf("ConfigFunc called %d times, want 1", callCount)
	}

	// Second call should return cached result
	cfg2, err := f.Config()
	if err != nil {
		t.Fatalf("Config() second call error: %v", err)
	}
	if cfg2 != expectedCfg {
		t.Error("Config() second call returned different config")
	}
	if callCount != 1 {
		t.Errorf("ConfigFunc called %d times after second call, want 1 (should be cached)", callCount)
	}
}

func TestEffectiveProfile_WithFlag(t *testing.T) {
	f := &Factory{
		IOStreams: DefaultIOStreams(),
		Profile:  "staging",
	}
	f.ConfigFunc = func() (*core.CliConfig, error) {
		return &core.CliConfig{ActiveProfile: "prod"}, nil
	}

	if got := f.EffectiveProfile(); got != "staging" {
		t.Errorf("EffectiveProfile() = %q, want %q", got, "staging")
	}
}

func TestEffectiveProfile_WithoutFlag_UsesConfig(t *testing.T) {
	f := &Factory{
		IOStreams: DefaultIOStreams(),
	}
	f.ConfigFunc = func() (*core.CliConfig, error) {
		return &core.CliConfig{ActiveProfile: "prod"}, nil
	}

	if got := f.EffectiveProfile(); got != "prod" {
		t.Errorf("EffectiveProfile() = %q, want %q", got, "prod")
	}
}

func TestEffectiveProfile_NoFlagNoConfig_DefaultsToDefault(t *testing.T) {
	f := &Factory{
		IOStreams: DefaultIOStreams(),
	}
	f.ConfigFunc = func() (*core.CliConfig, error) {
		return &core.CliConfig{}, nil
	}

	if got := f.EffectiveProfile(); got != "" {
		// When ActiveProfile is "" and no flag, Config returns "" which is what we get
		// Actually, looking at the code: if cfg.ActiveProfile is "", it returns ""
		// But the fallback is "default" only when Config() errors
		t.Errorf("EffectiveProfile() = %q, want empty string", got)
	}
}

func TestHttpClient_Default(t *testing.T) {
	f := NewFactory()

	client, err := f.HttpClient()
	if err != nil {
		t.Fatalf("HttpClient() error: %v", err)
	}
	if client == nil {
		t.Fatal("HttpClient() returned nil")
	}
}
