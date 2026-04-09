package cmdutil

import (
	"net/http"
	"sync"

	"github.com/madaima/qmai-cli/internal/core"
	"github.com/madaima/qmai-cli/internal/keychain"
)

// Factory provides lazy-initialized dependencies for commands
type Factory struct {
	IOStreams *IOStreams
	Keychain keychain.Keychain

	// Lazy-initialized
	configOnce sync.Once
	config     *core.CliConfig
	configErr  error

	httpClientOnce sync.Once
	httpClient     *http.Client
	httpClientErr  error

	// Overridable for testing
	ConfigFunc     func() (*core.CliConfig, error)
	HttpClientFunc func() (*http.Client, error)

	// Global flags
	Format  string
	Profile string
	Debug   bool
}

// NewFactory creates a Factory with default implementations
func NewFactory() *Factory {
	f := &Factory{
		IOStreams: DefaultIOStreams(),
	}
	f.ConfigFunc = func() (*core.CliConfig, error) {
		return core.LoadConfig()
	}
	return f
}

// Config returns the CLI configuration (lazy loaded, cached)
func (f *Factory) Config() (*core.CliConfig, error) {
	f.configOnce.Do(func() {
		f.config, f.configErr = f.ConfigFunc()
	})
	return f.config, f.configErr
}

// HttpClient returns an HTTP client with auth transport (lazy loaded)
func (f *Factory) HttpClient() (*http.Client, error) {
	f.httpClientOnce.Do(func() {
		if f.HttpClientFunc != nil {
			f.httpClient, f.httpClientErr = f.HttpClientFunc()
			return
		}
		// Default: plain client, auth transport added after login
		f.httpClient = &http.Client{}
	})
	return f.httpClient, f.httpClientErr
}

// EffectiveFormat returns the output format from flag -> config -> default
func (f *Factory) EffectiveFormat() string {
	if f.Format != "" {
		return f.Format
	}
	if cfg, err := f.Config(); err == nil {
		return cfg.Format()
	}
	return "table"
}

// EffectiveProfile returns the active profile name
func (f *Factory) EffectiveProfile() string {
	if f.Profile != "" {
		return f.Profile
	}
	if cfg, err := f.Config(); err == nil {
		return cfg.ActiveProfile
	}
	return "default"
}
