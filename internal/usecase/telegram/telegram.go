package telegram

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"telegram-fuse/internal/entity"
	"telegram-fuse/pkg/config"
	"telegram-fuse/pkg/repository"
)

type Bot struct {
	api *tgbotapi.BotAPI
	db  repository.Repository
}

func NewBot(api *tgbotapi.BotAPI, db repository.Repository) *Bot {
	return &Bot{api: api, db: db}
}

func (b *Bot) SaveFile(parentId int, name string, data []byte) (int, error) {
	file := tgbotapi.FileBytes{
		Name:  name,
		Bytes: data,
	}

	doc := tgbotapi.NewDocument(config.TelegramCfg.ChatId, file)
	doc.Caption = strconv.FormatInt(int64(parentId), 10)

	message, err := b.api.Send(doc)
	if err != nil {
		return 0, fmt.Errorf("couldn't send message: %w", err)
	}

	e := entity.FilesystemEntity{
		ParentId:  parentId,
		Name:      name,
		Size:      message.Document.FileSize,
		MessageID: message.MessageID,
		FileID:    message.Document.FileID,
	}

	if _, err := b.db.SaveEntity(e); err != nil {
		return 0, fmt.Errorf("couldn't save entity: %w", err)
	}

	return 0, err
}

func (b *Bot) ReadFile(id int) ([]byte, error) {
	e, err := b.db.GetEntity(id)
	if err != nil {
		return nil, fmt.Errorf("couldn't get entity: %w", err)
	}

	fileURL, err := b.api.GetFileDirectURL(e.FileID)
	if err != nil {
		return nil, fmt.Errorf("couldn't get file URL: %w", err)
	}

	resp, err := http.Get(fileURL)
	if err != nil {
		return nil, fmt.Errorf("couldn't get file from url: %w", err)
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func (b *Bot) GetDirectoryChildren(id int) ([]entity.FilesystemEntity, error) {
	return b.db.GetDirectoryChildren(id)
}
