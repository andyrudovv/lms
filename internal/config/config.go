package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App struct {
		Port int `yaml:"port"`
	} `yaml:"app"`

	DB struct {
		DSN                string `yaml:"dsn"`
		MaxConns           int32  `yaml:"max_conns"`
		MinConns           int32  `yaml:"min_conns"`
		MaxConnIdleMinutes int    `yaml:"max_conn_idle_minutes"`
	} `yaml:"db"`

	JWT struct {
		Secret           string `yaml:"secret"`
		AccessTTLMinutes int    `yaml:"access_ttl_minutes"`
	} `yaml:"jwt"`

	Migrations struct {
		Dir    string `yaml:"dir"`
		AutoUp bool   `yaml:"auto_up"`
	} `yaml:"migrations"`
}

func Load(path string) (Config, error) {
	if path == "" {
		path = "config.yaml"
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return Config{}, err
	}

	if cfg.App.Port == 0 {
		return Config{}, errors.New("config: app.port is required")
	}
	if cfg.DB.DSN == "" {
		return Config{}, errors.New("config: db.dsn is required")
	}
	if cfg.JWT.Secret == "" {
		return Config{}, errors.New("config: jwt.secret is required")
	}
	if cfg.JWT.AccessTTLMinutes == 0 {
		cfg.JWT.AccessTTLMinutes = 60
	}
	if cfg.Migrations.Dir == "" {
		cfg.Migrations.Dir = "migrations"
	}

	if cfg.DB.MaxConns == 0 {
		cfg.DB.MaxConns = 10
	}
	if cfg.DB.MinConns == 0 {
		cfg.DB.MinConns = 1
	}
	if cfg.DB.MaxConnIdleMinutes == 0 {
		cfg.DB.MaxConnIdleMinutes = 5
	}

	return cfg, nil
}
