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

func (b *Bot) DeleteEntity(id int) error {
	e, err := b.db.GetEntity(id)
	if err != nil {
		return fmt.Errorf("couldn't get entity: %w", err)
	}

	deleteMessageConfig := tgbotapi.NewDeleteMessage(config.TelegramCfg.ChatId, e.MessageID)
	_, err = b.api.Request(deleteMessageConfig)
	if err != nil {
		return fmt.Errorf("couldn't delete message: %w", err)
	}

	return b.db.DeleteEntity(id)
}

func (b *Bot) SaveDirectory(parentId int, name string) (entity.FilesystemEntity, error) {
	msg := tgbotapi.NewMessage(config.TelegramCfg.ChatId, fmt.Sprintf(
		"%s %s",
		name,
		strconv.FormatInt(int64(parentId), 10),
	))
	message, err := b.api.Send(msg)
	if err != nil {
		return entity.FilesystemEntity{}, fmt.Errorf("couldn't send message: %w", err)
	}

	e := entity.NewDirectory(
		0,
		parentId,
		name,
		message.MessageID,
		message.Time(),
		message.Time(),
	)

	entityId, err := b.db.SaveEntity(e)
	if err != nil {
		return entity.FilesystemEntity{}, fmt.Errorf("couldn't save entity: %w", err)
	}

	e.Id = entityId
	return e, err
}

func (b *Bot) UpdateEntity(filesystemEntity entity.FilesystemEntity) (*entity.FilesystemEntity, error) {
	// If the entity is a directory, update the message text
	if filesystemEntity.IsDirectory() {
		msg := tgbotapi.NewEditMessageText(
			config.TelegramCfg.ChatId,
			filesystemEntity.MessageID,
			fmt.Sprintf("%s %s", filesystemEntity.Name, strconv.FormatInt(int64(filesystemEntity.ParentId), 10)),
		)
		_, err := b.api.Send(msg)
		if err != nil {
			return nil, fmt.Errorf("couldn't send message: %w", err)
		}

		return &filesystemEntity, nil
	}

	// If the entity is a file, create a new message with the file content and delete the old one
	content, err := b.ReadFile(filesystemEntity.Id)
	if err != nil {
		return nil, fmt.Errorf("couldn't read file: %w", err)
	}

	file := tgbotapi.FileBytes{
		Name:  filesystemEntity.Name,
		Bytes: content,
	}

	doc := tgbotapi.NewDocument(config.TelegramCfg.ChatId, file)
	doc.Caption = strconv.FormatInt(int64(filesystemEntity.ParentId), 10)

	message, err := b.api.Send(doc)
	if err != nil {
		return nil, fmt.Errorf("couldn't send message: %w", err)
	}

	deleteOldMessageConfig := tgbotapi.NewDeleteMessage(config.TelegramCfg.ChatId, filesystemEntity.MessageID)
	_, err = b.api.Request(deleteOldMessageConfig)
	if err != nil {
		return nil, fmt.Errorf("couldn't delete old message: %w", err)
	}

	filesystemEntity.Size = message.Document.FileSize
	filesystemEntity.MessageID = message.MessageID
	filesystemEntity.FileID = message.Document.FileID
	filesystemEntity.UpdatedAt = message.Time()

	return &filesystemEntity, nil
}

func (b *Bot) UpdateFile(id int, data []byte) (entity.FilesystemEntity, error) {
	e, err := b.db.GetEntity(id)
	if err != nil {
		return entity.FilesystemEntity{}, fmt.Errorf("couldn't get entity: %w", err)
	}

	previousMessageID := e.MessageID
	file := tgbotapi.FileBytes{
		Name:  e.Name,
		Bytes: data,
	}

	doc := tgbotapi.NewDocument(config.TelegramCfg.ChatId, file)
	doc.Caption = strconv.FormatInt(int64(e.ParentId), 10)

	message, err := b.api.Send(doc)
	if err != nil {
		return entity.FilesystemEntity{}, fmt.Errorf("couldn't send message: %w", err)
	}

	e.Size = message.Document.FileSize
	e.MessageID = message.MessageID
	e.FileID = message.Document.FileID
	e.UpdatedAt = message.Time()

	err = b.db.UpdateEntity(e)
	if err != nil {
		return entity.FilesystemEntity{}, fmt.Errorf("couldn't update entity: %w", err)
	}

	deleteMessageConfig := tgbotapi.NewDeleteMessage(config.TelegramCfg.ChatId, previousMessageID)
	_, err = b.api.Request(deleteMessageConfig)
	if err != nil {
		return entity.FilesystemEntity{}, fmt.Errorf("couldn't delete message: %w", err)
	}

	return e, err
}

func (b *Bot) GetEntityById(id int) (entity.FilesystemEntity, error) {
	return b.db.GetEntity(id)
}

func NewBot(api *tgbotapi.BotAPI, db repository.Repository) *Bot {
	return &Bot{api: api, db: db}
}

func (b *Bot) SaveFile(parentId int, name string, data []byte) (entity.FilesystemEntity, error) {
	file := tgbotapi.FileBytes{
		Name:  name,
		Bytes: data,
	}

	doc := tgbotapi.NewDocument(config.TelegramCfg.ChatId, file)
	doc.Caption = strconv.FormatInt(int64(parentId), 10)

	message, err := b.api.Send(doc)
	if err != nil {
		return entity.FilesystemEntity{}, fmt.Errorf("couldn't send message: %w", err)
	}

	e := entity.NewFile(
		0,
		parentId,
		name,
		0,
		message.MessageID,
		message.Document.FileID,
		message.Time(),
		message.Time(),
	)

	entityId, err := b.db.SaveEntity(e)
	if err != nil {
		return entity.FilesystemEntity{}, fmt.Errorf("couldn't save entity: %w", err)
	}

	e.Id = entityId
	return e, err
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
