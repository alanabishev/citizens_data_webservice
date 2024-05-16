// Package config provides functionality for loading configuration data.
package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config is the main configuration structure.
// It includes the environment, storage path, and HTTP server configuration.
type Config struct {
	Env         string `yaml:"env" env-default:"local"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
}

// HTTPServer is a structure for HTTP server configuration.
// It includes the address, timeout, idle timeout, user, and password.
type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
	User        string        `yaml:"user" env-required:"true"`
	Password    string        `yaml:"password" env-required:"true" env:"HTTP_SERVER_PASSWORD"`
}

// Load reads the configuration file specified by the CONFIG_PATH environment variable.
// It returns a pointer to a Config structure.
// If the CONFIG_PATH environment variable is not set or the file does not exist, it logs a fatal error.
// If the configuration file cannot be read, it also logs a fatal error.
func Load() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}

	return &cfg
}
