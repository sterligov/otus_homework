package config

import (
	"fmt"
	"os"
	"time"

	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/logger"
	yaml "gopkg.in/yaml.v2"
)

const (
	SQLStorage      = "sql"
	InMemoryStorage = "in_memory"
)

type Config struct {
	HTTP struct {
		Addr           string        `yaml:"addr"`
		ReadTimeout    time.Duration `yaml:"read_timeout"`
		WriteTimeout   time.Duration `yaml:"write_timeout"`
		HandlerTimeout time.Duration `yaml:"handler_timeout"`
	} `yaml:"server"`

	Database struct {
		Addr   string `yaml:"connection_addr"`
		Driver string `yaml:"driver"`
	} `yaml:"database"`

	Logger struct {
		Path  string `yaml:"path"`
		Level string `yaml:"level"`
	} `yaml:"logger"`

	StorageType string `yaml:"storage_type"`
}

func New(cfgFilename string) (*Config, error) {
	f, err := os.Open(cfgFilename)
	if err != nil {
		return nil, fmt.Errorf("open config file failed: %w", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			logger.Warnf("config file close: %s", err)
		}
	}()

	cfg := &Config{}

	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(cfg); err != nil {
		return nil, fmt.Errorf("decode config file failed: %w", err)
	}

	return cfg, err
}
