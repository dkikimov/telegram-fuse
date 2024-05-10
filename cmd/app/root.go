package app

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/spf13/cobra"

	"telegram-fuse/internal/fuse"
	"telegram-fuse/pkg/config"
)

const (
	defaultTelegramConfig = "/etc/telegram-fuse/telegram.yaml"
)

func init() {
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", defaultTelegramConfig, "path to config")
	RootCmd.Flags().StringVar(&mountPath, "mount", "", "path to mount")
	if err := RootCmd.MarkFlagRequired("mount"); err != nil {
		panic(fmt.Sprintf("couldn't mark flag as required: %s", err))
	}
}

var cfgFile string
var mountPath string

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

		server, err := fs.Mount(mountPath, &fuse.TgNode{}, nil)
		if err != nil {
			return fmt.Errorf("couldn't mount fuse: %w", err)
		}

		slog.Info("mounted fuse")

		server.Serve()
		server.Wait()

		return nil
	},
}
