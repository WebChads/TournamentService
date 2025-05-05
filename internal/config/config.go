package config

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"

	slogerr "github.com/WebChads/TournamentService/internal/pkg/logger"
	"github.com/ilyakaznacheev/cleanenv"
)

type ServerConfig struct {
	LogLevel    string `yaml:"log_level"`
	Address     string `yaml:"address"`
	DatabaseURL string `yaml:"database_url"`
	SecretKey   string `yaml:"secret_key"`
}

func NewServerConfig() *ServerConfig {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		root, err := FindModuleRoot(".")
		if err != nil {
			slog.Error("failed to find config file", slogerr.Error(err))
			return nil
		}

		configPath = root + "/configs/local.yaml"
	}

	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		slog.Error("config file does not exists: "+configPath, slogerr.Error(err))
		return nil
	}

	var cfg ServerConfig

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		slog.Error("cannot read config", slogerr.Error(err))
		return nil
	}

	return &cfg
}

func FindModuleRoot(dir string) (string, error) {
	for {
		if dir == "" || dir == "/" {
			return "", errors.New("invalid working directory name")
		}

		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		dir = filepath.Dir(dir)
	}
}
