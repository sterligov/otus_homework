package config

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
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
	} `yaml:"http"`

	GRPC struct {
		Addr string `yaml:"addr"`
	} `yaml:"grpc"`

	StorageType string `yaml:"storage_type"`

	Database struct {
		Addr   string `yaml:"connection_addr"`
		Driver string `yaml:"driver"`
	}

	Logger struct {
		Path  string `yaml:"path"`
		Level string `yaml:"level"`
	}

	AMQP struct {
		ConnectionAddr      string        `yaml:"connection_addr"`
		QueueName           string        `yaml:"queue_name"`
		MaxReconnectRetries int           `yaml:"max_reconnect_retries"`
		ReconnectInterval   time.Duration `yaml:"reconnect_interval"`
		HandlersNumber      int           `yaml:"handlers_number"`
	}

	EventScanFreq time.Duration `yaml:"event_scan_frequency"`
}

func New(cfgFilename string) (*Config, error) {
	f, err := os.Open(cfgFilename)
	if err != nil {
		return nil, fmt.Errorf("open config file failed: %w", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			logrus.Warnf("config file close: %s", err)
		}
	}()

	cfg := &Config{}

	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(cfg); err != nil {
		return nil, fmt.Errorf("decode config file failed: %w", err)
	}

	return cfg, err
}
