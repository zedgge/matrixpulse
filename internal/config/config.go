package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Symbols     []string    `yaml:"symbols"`
	WindowSize  int         `yaml:"window_size"`
	UpdateHz    int         `yaml:"update_hz"`
	Alerts      Alerts      `yaml:"alerts"`
	Persistence Persistence `yaml:"persistence"`
	WebSocket   WebSocket   `yaml:"websocket"`
	REST        REST        `yaml:"rest"`
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

type WebSocket struct {
	Enabled bool `yaml:"enabled"`
	Port    int  `yaml:"port"`
}

type REST struct {
	Enabled bool `yaml:"enabled"`
	Port    int  `yaml:"port"`
}

type Dashboard struct {
	RefreshMs int  `yaml:"refresh_ms"`
	Enabled   bool `yaml:"enabled"`
}

func Load() *Config {
	cfg := &Config{
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
		WebSocket: WebSocket{Enabled: true, Port: 8080},
		REST:      REST{Enabled: true, Port: 8081},
		Dashboard: Dashboard{RefreshMs: 200, Enabled: true},
	}

	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return cfg
	}

	yaml.Unmarshal(data, cfg)
	return cfg
}
