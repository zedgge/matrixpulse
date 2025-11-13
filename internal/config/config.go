package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Symbols     []string    `yaml:"symbols"`
	WindowSize  int         `yaml:"window_size"`
	UpdateHz    int         `yaml:"update_hz"`
	Alerts      Alerts      `yaml:"alerts"`
	Persistence Persistence `yaml:"persistence"`
	Dashboard   Dashboard   `yaml:"dashboard"`
}

type Alerts struct {
	Correlation float64 `yaml:"correlation_threshold"`
	Eigenvalue  float64 `yaml:"eigenvalue_threshold"`
	Volatility  float64 `yaml:"volatility_threshold"`
}

type Persistence struct {
	Enabled  bool   `yaml:"enabled"`
	Path     string `yaml:"path"`
	Interval int    `yaml:"interval_seconds"`
}

type Dashboard struct {
	RefreshMs int  `yaml:"refresh_ms"`
	Enabled   bool `yaml:"enabled"`
}

// Load reads configuration from config.yaml with fallback to defaults
func Load() (*Config, error) {
	cfg := defaultConfig()

	data, err := os.ReadFile("config.yaml")
	if err != nil {
		if os.IsNotExist(err) {
			// Config file doesn't exist, use defaults
			return cfg, nil
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// defaultConfig returns sensible defaults
func defaultConfig() *Config {
	return &Config{
		Symbols:    []string{"AAPL", "GOOGL", "MSFT", "AMZN", "TSLA", "META"},
		WindowSize: 120,
		UpdateHz:   40,
		Alerts: Alerts{
			Correlation: 0.82,
			Eigenvalue:  2.8,
			Volatility:  0.04,
		},
		Persistence: Persistence{
			Enabled:  true,
			Path:     "matrixpulse_state.json",
			Interval: 60,
		},
		Dashboard: Dashboard{
			RefreshMs: 200,
			Enabled:   true,
		},
	}
}

// Validate checks configuration for errors
func (c *Config) Validate() error {
	if len(c.Symbols) == 0 {
		return fmt.Errorf("must specify at least one symbol")
	}

	if len(c.Symbols) > 100 {
		return fmt.Errorf("too many symbols (max 100, got %d)", len(c.Symbols))
	}

	if c.WindowSize < 10 {
		return fmt.Errorf("window_size too small (min 10, got %d)", c.WindowSize)
	}

	if c.WindowSize > 10000 {
		return fmt.Errorf("window_size too large (max 10000, got %d)", c.WindowSize)
	}

	if c.UpdateHz < 1 || c.UpdateHz > 1000 {
		return fmt.Errorf("update_hz out of range (1-1000, got %d)", c.UpdateHz)
	}

	if c.Alerts.Correlation < 0 || c.Alerts.Correlation > 1 {
		return fmt.Errorf("correlation_threshold must be 0-1 (got %.2f)", c.Alerts.Correlation)
	}

	if c.Alerts.Eigenvalue < 0 {
		return fmt.Errorf("eigenvalue_threshold must be positive (got %.2f)", c.Alerts.Eigenvalue)
	}

	if c.Persistence.Interval < 1 {
		return fmt.Errorf("persistence interval must be positive (got %d)", c.Persistence.Interval)
	}

	if c.Dashboard.RefreshMs < 50 || c.Dashboard.RefreshMs > 5000 {
		return fmt.Errorf("dashboard refresh_ms out of range (50-5000, got %d)", c.Dashboard.RefreshMs)
	}

	return nil
}
