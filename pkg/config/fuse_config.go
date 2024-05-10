package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// FuseCfg is the global configuration for the FUSE.
var FuseCfg *FuseConfig

// FuseConfig represents the configuration for the FUSE.
type FuseConfig struct {
	Debug     bool   `yaml:"debug"`
	MountPath string `yaml:"mount_path"`
}

// InitFuseConfigFromFile reads config file and initializes FUSE config.
func InitFuseConfigFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if err := initFuseFromString(data); err != nil {
		return err
	}

	return nil
}

// initFuseFromString initializes fuse config from a slice of bytes.
func initFuseFromString(data []byte) error {
	var cfg FuseConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("unable to unmarshal fuse config file: %w", err)
	}

	FuseCfg = &cfg

	return nil
}
