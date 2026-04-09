package core

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	DefaultBaseURL = "https://openapi.qmai.cn/"
	ConfigDir      = ".config/qmai"
	ConfigFile     = "config.yaml"
)

type Profile struct {
	Name      string `yaml:"name"`
	ShopCode  string `yaml:"shop_code,omitempty"`  // 门店编码
	BaseURL   string `yaml:"base_url,omitempty"`
	OpenId    string `yaml:"open_id,omitempty"`    // 开放平台 openId
	GrantCode string `yaml:"grant_code,omitempty"` // 开放平台 grantCode
}

type CliConfig struct {
	ActiveProfile string              `yaml:"active_profile"`
	DefaultFormat string              `yaml:"default_format,omitempty"` // json, table, csv
	Debug         bool                `yaml:"debug,omitempty"`
	Profiles      map[string]*Profile `yaml:"profiles,omitempty"`
	path          string              // file path, not serialized
}

// ConfigPath returns the full path to the config file
func ConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine home directory: %w", err)
	}
	return filepath.Join(home, ConfigDir, ConfigFile), nil
}

// LoadConfig loads config from the default path
func LoadConfig() (*CliConfig, error) {
	path, err := ConfigPath()
	if err != nil {
		return nil, err
	}
	return LoadConfigFrom(path)
}

// LoadConfigFrom loads config from a specific path
func LoadConfigFrom(path string) (*CliConfig, error) {
	cfg := &CliConfig{
		Profiles: make(map[string]*Profile),
		path:     path,
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil // return empty config, not an error
		}
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	cfg.path = path
	if cfg.Profiles == nil {
		cfg.Profiles = make(map[string]*Profile)
	}
	return cfg, nil
}

// Save writes config to disk
func (c *CliConfig) Save() error {
	if c.path == "" {
		var err error
		c.path, err = ConfigPath()
		if err != nil {
			return err
		}
	}
	dir := filepath.Dir(c.path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	return os.WriteFile(c.path, data, 0o644)
}

// ActiveProfileConfig returns the currently active profile, or nil
func (c *CliConfig) ActiveProfileConfig() *Profile {
	if c.ActiveProfile == "" {
		return nil
	}
	return c.Profiles[c.ActiveProfile]
}

// BaseURL returns the API base URL for the active profile
func (c *CliConfig) BaseURL() string {
	if p := c.ActiveProfileConfig(); p != nil && p.BaseURL != "" {
		return p.BaseURL
	}
	return DefaultBaseURL
}

// Format returns the effective output format
func (c *CliConfig) Format() string {
	if c.DefaultFormat != "" {
		return c.DefaultFormat
	}
	return "table"
}

// SetProfile adds or updates a profile
func (c *CliConfig) SetProfile(name string, p *Profile) {
	p.Name = name
	c.Profiles[name] = p
}

// Path returns the config file path
func (c *CliConfig) Path() string {
	return c.path
}
