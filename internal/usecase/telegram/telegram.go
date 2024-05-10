package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"telegram-fuse/internal/repository/sqlite"
	"telegram-fuse/pkg/config"
)

type Bot struct {
	api *tgbotapi.BotAPI
	db  *sqlite.Database
}

func NewBot(api *tgbotapi.BotAPI, db *sqlite.Database) *Bot {
	return &Bot{api: api, db: db}
}

func (b *Bot) SaveFile(path string, name string, data []byte) error {
	file := tgbotapi.FileBytes{
		Name:  name,
		Bytes: data,
	}

	doc := tgbotapi.NewDocument(config.TelegramCfg.ChatId, file)
	doc.Caption = path

	_, err := b.api.Send(doc)

	return err
}
