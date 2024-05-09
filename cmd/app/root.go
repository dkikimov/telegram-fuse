package app

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"telegram-fuse/pkg/config"
)

const (
	defaultTelegramConfig = "/etc/telegram-fuse/telegram.yaml"
)

func init() {
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", defaultTelegramConfig, "path to config")
}

var cfgFile string

var RootCmd = &cobra.Command{
	Use:   "tgfuse",
	Short: "tgfuse is a fuse implementation for telegram bot",
	RunE: func(_ *cobra.Command, _ []string) error {
		if _, err := os.Stat(cfgFile); err != nil {
			return fmt.Errorf("unable to read config file %s: %w", cfgFile, err)
		}

		if err := config.InitTelegramConfigFromFile(cfgFile); err != nil {
			return fmt.Errorf("unable to init config: %w", err)
		}

		return nil
	},
}
