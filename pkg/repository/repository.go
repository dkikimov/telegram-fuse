package repository

import "telegram-fuse/internal/entity"

type Repository interface {
	Close() error
	SaveEntity(filesystemEntity entity.FilesystemEntity) error
	GetDirectoryChildren(filesystemEntity entity.FilesystemEntity) ([]entity.FilesystemEntity, error)
	RenameEntity(filesystemEntity entity.FilesystemEntity) error
	RemoveEntity(filesystemEntity entity.FilesystemEntity) error
}