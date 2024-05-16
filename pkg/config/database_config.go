package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var DatabaseCfg *DatabaseConfig

type DatabaseConfig struct {
	Path string `yaml:"path"`
}

func InitDatabaseConfigFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if err := initDatabaseFromString(data); err != nil {
		return err
	}

	return nil
}

func initDatabaseFromString(data []byte) error {
	var cfg DatabaseConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("unable to unmarshal fuse config file: %w", err)
	}

	DatabaseCfg = &cfg

	return nil
}
