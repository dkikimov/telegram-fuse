package app

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/spf13/cobra"

	"telegram-fuse/internal/repository/sqlite"
	"telegram-fuse/internal/tgfuse"
	"telegram-fuse/internal/usecase/telegram"
	"telegram-fuse/pkg/config"
)

const (
	defaultTelegramConfig = "/etc/telegram-fuse/telegram.yaml"
	defaultFuseConfig     = "/etc/telegram-fuse/fuse.yaml"
)

func init() {
	RootCmd.PersistentFlags().StringVar(&telegramCfgFile, "telegram-config", defaultTelegramConfig, "path to telegram config")
	RootCmd.PersistentFlags().StringVar(&fuseCfgFile, "fuse-config", defaultFuseConfig, "path to fuse config")
}

var telegramCfgFile string
var fuseCfgFile string

var RootCmd = &cobra.Command{
	Use:   "tgfuse",
	Short: "tgfuse is a fuse implementation for telegram bot",
	RunE: func(_ *cobra.Command, _ []string) error {
		// Init cfg
		if _, err := os.Stat(telegramCfgFile); err != nil {
			return fmt.Errorf("unable to read config file %s: %w", telegramCfgFile, err)
		}

		if _, err := os.Stat(fuseCfgFile); err != nil {
			return fmt.Errorf("unable to read config file %s: %w", fuseCfgFile, err)
		}

		if err := config.InitTelegramConfigFromFile(telegramCfgFile); err != nil {
			return fmt.Errorf("unable to init telegram config: %w", err)
		}

		if err := config.InitFuseConfigFromFile(fuseCfgFile); err != nil {
			return fmt.Errorf("unable to init fuse config: %w", err)
		}

		// Init database
		db, err := sqlite.NewDatabase()
		if err != nil {
			return fmt.Errorf("couldn't initialize database with error: %w", err)
		}

		// Init telegram bot
		api, err := tgbotapi.NewBotAPI(config.TelegramCfg.Token)
		if err != nil {
			return fmt.Errorf("couldn't initialize bot with error: %w", err)
		}

		bot := telegram.NewBot(api, db)

		// Init fuse
		timeout := time.Second
		server, err := fs.Mount(config.FuseCfg.MountPath, tgfuse.NewNode(bot), &fs.Options{
			MountOptions: fuse.MountOptions{
				Debug: config.FuseCfg.Debug,
			},
			EntryTimeout: &timeout,
			AttrTimeout:  &timeout,
		})

		if err != nil {
			return fmt.Errorf("couldn't mount fuse: %w", err)
		}

		slog.Info("mounted fuse")

		go server.Serve()
		server.Wait()

		// TODO: catch signal
		if err := server.Unmount(); err != nil {
			return fmt.Errorf("couldn't unmount fuse: %s", err)
		}

		slog.Info("unmounted fuse")
		return nil
	},
}
