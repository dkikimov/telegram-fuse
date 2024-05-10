package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// TelegramCfg is the global configuration for the Telegram bot.
var TelegramCfg *TelegramConfig

// TelegramConfig represents the configuration for the Telegram bot.
type TelegramConfig struct {
	Token  string `yaml:"token"`
	ChatId int64  `yaml:"chat_id"`
}

// InitTelegramConfigFromFile reads config file and initializes telegram config.
func InitTelegramConfigFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if err := initTelegramFromString(data); err != nil {
		return err
	}

	return nil
}

// initTelegramFromString initializes telegram config from a slice of bytes.
func initTelegramFromString(data []byte) error {
	var cfg TelegramConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("unable to unmarshal config file: %w", err)
	}

	TelegramCfg = &cfg

	return nil
}
