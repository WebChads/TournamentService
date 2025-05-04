package config

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"

	slogerr "github.com/WebChads/TournamentService/internal/pkg/logger"
	"github.com/ilyakaznacheev/cleanenv"
)

type ServerEnv struct {
	Address     string `yaml:"address"`
	DatabaseURL string `yaml:"database_url"`
}

type ServerConfig struct {
	LogLevel  string    `yaml:"log_level"`
	LocalEnv  ServerEnv `yaml:"local"`
	DockerEnv ServerEnv `yaml:"docker"`
}

const (
	serverEnvLocal  = "local"
	serverEnvDocker = "docker"
)

func NewServerConfig() (*ServerConfig, *ServerEnv) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		root, err := FindModuleRoot(".")
		if err != nil {
			slog.Error("failed to find config file", slogerr.Error(err))
			return nil, nil
		}

		configPath = root + "/configs/local.yaml"
	}

	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		slog.Error("config file does not exists: "+configPath, slogerr.Error(err))
		return nil, nil
	}

	var cfg ServerConfig

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		slog.Error("cannot read config", slogerr.Error(err))
		return nil, nil
	}

	env := os.Getenv("SERVER_ENV")
	if env == "" {
		// set default server environment to local env
		env = serverEnvLocal
	}

	var serverEnv ServerEnv
	switch env {
	case serverEnvLocal:
		serverEnv = ServerEnv{
			Address:     cfg.LocalEnv.Address,
			DatabaseURL: cfg.LocalEnv.DatabaseURL,
		}
	case serverEnvDocker:
		serverEnv = ServerEnv{
			Address:     cfg.DockerEnv.Address,
			DatabaseURL: cfg.DockerEnv.DatabaseURL,
		}
	}

	return &cfg, &serverEnv
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
