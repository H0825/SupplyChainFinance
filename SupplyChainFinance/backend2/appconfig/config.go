package appconfig

import (
	"fmt"
	"net/url"
	"os"
	"sync"

	toml "github.com/pelletier/go-toml/v2"
)

// Settings stores backend runtime configuration.
type Settings struct {
	App struct {
		Listen       string `toml:"listen"`
		LogFile      string `toml:"log_file"`
		FrontendDist string `toml:"frontend_dist"`
		LegacyIndex  string `toml:"legacy_index"`
		LegacyAssets string `toml:"legacy_assets"`
	} `toml:"App"`
	Database struct {
		User      string `toml:"user"`
		Password  string `toml:"password"`
		Host      string `toml:"host"`
		Port      int    `toml:"port"`
		Name      string `toml:"name"`
		Charset   string `toml:"charset"`
		ParseTime bool   `toml:"parse_time"`
		Loc       string `toml:"loc"`
	} `toml:"Database"`
	Contract struct {
		Address string `toml:"address"`
	} `toml:"Contract"`
	AI struct {
		BaseURL     string  `toml:"base_url"`
		APIKey      string  `toml:"api_key"`
		Model       string  `toml:"model"`
		TimeoutSec  int     `toml:"timeout_sec"`
		Temperature float64 `toml:"temperature"`
	} `toml:"AI"`
}

var (
	once    sync.Once
	cached  *Settings
	loadErr error
)

func applyDefaults(cfg *Settings) {
	if cfg.App.Listen == "" {
		cfg.App.Listen = ":8888"
	}
	if cfg.App.LogFile == "" {
		cfg.App.LogFile = "gin-runner.log"
	}
	if cfg.App.FrontendDist == "" {
		cfg.App.FrontendDist = "../frontend/dist"
	}
	if cfg.App.LegacyIndex == "" {
		cfg.App.LegacyIndex = "./web/index.html"
	}
	if cfg.App.LegacyAssets == "" {
		cfg.App.LegacyAssets = "./web/assets"
	}
	if cfg.Database.Host == "" {
		cfg.Database.Host = "localhost"
	}
	if cfg.Database.Port == 0 {
		cfg.Database.Port = 3306
	}
	if cfg.Database.Charset == "" {
		cfg.Database.Charset = "utf8mb4"
	}
	if cfg.Database.Loc == "" {
		cfg.Database.Loc = "Local"
	}
	if cfg.Database.Name == "" {
		cfg.Database.Name = "finance"
	}
	if cfg.AI.BaseURL == "" {
		cfg.AI.BaseURL = "https://ark.cn-beijing.volces.com/api/v3"
	}
	if cfg.AI.Model == "" {
		cfg.AI.Model = "doubao-1-5-pro-32k-250115"
	}
	if cfg.AI.TimeoutSec <= 0 {
		cfg.AI.TimeoutSec = 30
	}
	if cfg.AI.Temperature == 0 {
		cfg.AI.Temperature = 0.2
	}
	if cfg.AI.APIKey == "" {
		cfg.AI.APIKey = os.Getenv("AI_API_KEY")
	}
}

// Load loads config.toml once and returns cached settings.
func Load(path string) (*Settings, error) {
	once.Do(func() {
		cfg := &Settings{}
		b, err := os.ReadFile(path)
		if err != nil {
			loadErr = fmt.Errorf("read config failed: %w", err)
			return
		}
		if err = toml.Unmarshal(b, cfg); err != nil {
			loadErr = fmt.Errorf("parse config failed: %w", err)
			return
		}
		applyDefaults(cfg)
		cached = cfg
	})
	return cached, loadErr
}

// DSN builds mysql dsn from split database fields.
func (s *Settings) DSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		s.Database.User,
		s.Database.Password,
		s.Database.Host,
		s.Database.Port,
		s.Database.Name,
		s.Database.Charset,
		s.Database.ParseTime,
		url.QueryEscape(s.Database.Loc),
	)
}
